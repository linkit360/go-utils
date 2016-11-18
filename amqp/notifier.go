package rabbit

// this package for holding connection to rabbit.
// it has 2 channels, reading and pending
// it reconnects is the connection has lost
// metrics avialable, do not forget to add handler for /var/debug

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	amqp_driver "github.com/streadway/amqp"

	m "github.com/vostrok/utils/metrics"
)

type Notifier struct {
	url            string
	reconnectDelay int
	done           chan error
	conn           *amqp_driver.Connection
	channel        *amqp_driver.Channel
	m              NotifierMetrics
	publishCh      chan AMQPMessage
}
type ConnectionConfig struct {
	User string `yaml:"user" default:"linkit"`
	Pass string `yaml:"pass" default:"dg-U_oHhy7-"`
	Host string `yaml:"host" default:"35.154.8.158"`
	Port string `yaml:"port" default:"5672"`
}
type NotifierConfig struct {
	Conn           ConnectionConfig `yaml:"conn"`
	ReconnectDelay int              `default:"10" yaml:"reconnect_delay"`
	ChanCapacity   int64            `default:"1000000" yaml:"chan_capacity"`
}

func NewNotifier(c NotifierConfig) *Notifier {
	connectionUrl := fmt.Sprintf("amqp://%s:%s@%s:%s", c.Conn.User, c.Conn.Pass, c.Conn.Host, c.Conn.Port)
	notifier := &Notifier{
		url:            connectionUrl,
		reconnectDelay: c.ReconnectDelay,
		done:           make(chan error),
		conn:           nil,
		channel:        nil,
		m:              initNotifierMetrics(),
		publishCh:      make(chan AMQPMessage, c.ChanCapacity),
	}
	if err := notifier.connect(); err != nil {
		log.Fatal("Connect error", err.Error())
	}
	go func() {
		notifier.publisher()
	}()

	return notifier
}

func (n *Notifier) Publish(msg AMQPMessage) {
	n.publishCh <- msg
}

func (n *Notifier) GetQueueSize(queue string) (int, error) {
	queueInfo, err := n.channel.QueueInspect(queue)
	if err != nil {
		err = fmt.Errorf("channel.QueueInspect: %s", err.Error())
		log.WithFields(log.Fields{
			"queue": queue,
			"error": err.Error(),
		}).Error("rbmq consumer: cannot inspect queue")
		return 0, err
	}
	return queueInfo.Messages, nil
}

func (n *Notifier) connect() error {
	var err error

	log.WithField("url", n.url).Debug("dialing")
	n.conn, err = amqp_driver.Dial(n.url)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		// Waits here for the channel to be closed
		log.Info("closing: %s", <-n.conn.NotifyClose(make(chan *amqp_driver.Error)))
		// Let Handle know it's not time to reconnect
		n.done <- errors.New("Channel Closed")
	}()

	log.Info("rbmq notifier: got connection, getting channel...")
	n.channel, err = n.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	log.Info("rbmq notifier: got channel")
	return nil
}

func (n *Notifier) reConnect() error {
	log.WithField("reconnectDelay", n.reconnectDelay).Info("consumer reconnects...")
	time.Sleep(time.Duration(n.reconnectDelay) * time.Second)

	if err := n.connect(); err != nil {
		n.m.Connected.Set(0)
		n.m.ReconnectCount.Inc()

		log.WithField("error", err.Error()).Error("could not reconnect")
		return fmt.Errorf("Connect: %s", err.Error())
	}
	n.m.Connected.Set(1)
	n.m.ReconnectCount.Set(0)
	return nil
}

type EventNotify struct {
	EventName string      `json:"event_name,omitempty"`
	EventData interface{} `json:"event_data,omitempty"`
}
type AMQPMessage struct {
	QueueName string
	Body      []byte
}

type session struct {
	*amqp_driver.Connection
	*amqp_driver.Channel
}

func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

func (n *Notifier) publisher() {

	var (
		running bool
		reading = n.publishCh
		pending = make(chan AMQPMessage, cap(n.publishCh))
	)

	go func() {
		for range time.Tick(1 * time.Second) {
			n.m.PendingBuffer.Set(float64(len(pending)))
			n.m.ReadingBuffer.Set(float64(len(reading)))
		}
	}()

	go func() {
		for {
			var msg AMQPMessage
			msg, running = <-reading
			if !running {
				log.WithField("rbmq", "!running").Info("rbmq: publisher")
				return
			}

			if pending <- msg; len(pending) > 0 {
				log.WithField("rbmq_pending", len(pending)).Info("rbmq: publisher")
			}
		}
	}()

	for {
		var msg AMQPMessage
		select {

		case <-n.done:
			err := n.reConnect()
			if err != nil {
				// Very likely chance of failing
				// should not cause worker to terminate
				log.WithField("error", err.Error()).Error("rbmq notifier reconnect failed")
			} else {
				log.Info("rbmq notifier: reconnected")
			}

		case msg = <-pending:

			q, err := n.channel.QueueDeclare(
				msg.QueueName, // name
				false,         // durable
				false,         // delete when unused
				false,         // exclusive
				false,         // no-wait
				nil,           // arguments
			)

			if err != nil {
				pending <- msg
				log.WithField("error", err.Error()).Error("rbmq: Channel.QueueDeclare")
				n.m.PublishErrs.Inc()
				n.conn.Close()
				break
			}

			err = n.channel.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp_driver.Publishing{
					ContentType: "text/plain",
					Body:        msg.Body,
				})

			if err != nil {
				pending <- msg
				log.WithField("error", err.Error()).Error("rbmq: publish")
				n.m.PublishErrs.Inc()
				n.conn.Close()
				break
			}
			log.WithFields(log.Fields{"queue": msg.QueueName, "body": string(msg.Body)}).Info("rbmq: publish")
		}
	}
}

type NotifierMetrics struct {
	SessionRequests m.Gauge
	PublishErrs     m.Gauge

	ReconnectCount prometheus.Gauge
	Connected      prometheus.Gauge
	PendingBuffer  prometheus.Gauge
	ReadingBuffer  prometheus.Gauge
}

func newGaugePublisher(name, help string) m.Gauge {
	return m.NewGauge("rbmq", "notifier", name, "rbmq "+help)
}
func initNotifierMetrics() NotifierMetrics {
	metrics := NotifierMetrics{
		SessionRequests: newGaugePublisher("reconnects_count", "publisher reconnect count"),
		PublishErrs:     newGaugePublisher("publish_errors", "publish errors"),

		Connected:      m.PrometheusGauge("rbmq", "notifier", "connected", "publisher connection status"),
		ReconnectCount: m.PrometheusGauge("rbmq", "notifier", "reconnect_count", "publisher connection attempts count"),
		PendingBuffer:  m.PrometheusGauge("rbmq", "notifier", "buffer_pending_gauge_size", "publisher pending buffer size"),
		ReadingBuffer:  m.PrometheusGauge("rbmq", "notifier", "buffer_reading_gauge_size", "publisher reading buffer size"),
	}
	return metrics
}
