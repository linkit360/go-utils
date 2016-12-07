package amqp

// this package for holding connection to rabbit.
// it has 2 channels, reading and pending
// it reconnects is the connection has lost
// metrics avialable, do not forget to add handler for /var/debug

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	amqp_driver "github.com/streadway/amqp"

	m "github.com/vostrok/utils/metrics"
)

type Notifier struct {
	url            string
	conf           NotifierConfig
	reconnectDelay int
	stop           bool
	done           chan error
	conn           *amqp_driver.Connection
	channel        *amqp_driver.Channel
	m              NotifierMetrics
	publishCh      chan AMQPMessage
	signalsCh      chan os.Signal
	FinishCh       chan bool
}
type ConnectionConfig struct {
	User string `yaml:"user" default:"linkit"`
	Pass string `yaml:"pass" default:"dg-U_oHhy7-"`
	Host string `yaml:"host" default:"localhost"`
	Port string `yaml:"port" default:"5672"`
}
type NotifierConfig struct {
	Conn           ConnectionConfig `yaml:"conn"`
	ReconnectDelay int              `default:"10" yaml:"reconnect_delay"`
	ChanCapacity   int64            `default:"1000000" yaml:"chan_capacity"`
	BufferPath     string           `yaml:"pending_buffer_path"`
}

func NewNotifier(c NotifierConfig) *Notifier {
	connectionUrl := fmt.Sprintf("amqp://%s:%s@%s:%s", c.Conn.User, c.Conn.Pass, c.Conn.Host, c.Conn.Port)
	notifier := &Notifier{
		url:            connectionUrl,
		conf:           c,
		reconnectDelay: c.ReconnectDelay,
		done:           make(chan error),
		conn:           nil,
		channel:        nil,
		m:              initNotifierMetrics(),
		publishCh:      make(chan AMQPMessage, c.ChanCapacity),
		signalsCh:      make(chan os.Signal, 1),
		FinishCh:       make(chan bool, 1),
	}
	if notifier.conf.BufferPath == "" {
		log.Fatal("No buffer path (pending_buffer_path)")
	}
	if err := notifier.connect(); err != nil {
		log.Fatal("Connect error ", err.Error())
	}

	notifier.restoreState()

	go func() {
		notifier.publisher()
	}()

	signal.Notify(notifier.signalsCh, syscall.SIGINT, syscall.SIGTERM)

	return notifier
}

func (n *Notifier) Publish(msg AMQPMessage) {
	n.publishCh <- msg
}

func (n *Notifier) restoreState() {

	fh, err := os.Open(n.conf.BufferPath)
	if err != nil {
		log.WithField("error", err.Error()).Info("cannot open pending buffer file")
		return
	}
	bufferBytes := bytes.NewBuffer(nil)
	_, err = io.Copy(bufferBytes, fh)
	if err != nil {
		log.WithField("error", err.Error()).Error("rbmq notifier cannot copy from buffer file")
		return
	}
	if err := fh.Close(); err != nil {
		log.WithField("error", err.Error()).Error("rbmq notifier cannot close buffer fh")
		return
	}
	var buf []AMQPMessage
	if err := json.Unmarshal(bufferBytes.Bytes(), buf); err != nil {
		log.WithField("error", err.Error()).Error("rbmq notifier cannot unmarshal")
		return
	}
	log.WithField("count", len(buf)).Debug("rbmq notifier restore state")
	for _, msg := range buf {
		n.Publish(msg)
	}
}

type Buffer struct {
	Reading chan AMQPMessage `json:"reading"`
	Pending chan AMQPMessage `json:"pending"`
}

func (n *Notifier) GetQueueSize(queue string) (int, error) {
getQueueSize:
	queueInfo, err := n.channel.QueueInspect(queue)
	if err != nil {
		if err == amqp_driver.ErrClosed {
			log.Info("rbmq notifier try to reconnect")
			if err := n.reConnect(); err != nil {
				log.WithField("error", err.Error()).Error("rbmq notifier reconnect failed")
				return 0, err
			} else {

				goto getQueueSize
			}
		}
		err = fmt.Errorf("channel.QueueInspect: %s", err.Error())
		log.WithFields(log.Fields{
			"queue": queue,
			"error": err.Error(),
		}).Error("rbmq notifier: cannot inspect queue")
		return 0, err
	}
	return queueInfo.Messages, nil
}

func (n *Notifier) connect() error {
	var err error

	log.WithField("url", n.url).Debug("rbmq notifier dialing...")
	n.conn, err = amqp_driver.Dial(n.url)
	if err != nil {
		return fmt.Errorf("amqp_driver.Dial: %s", err)
	}

	go func() {
		// Waits here for the channel to be closed
		log.Info("rbmq notifier closing: ", <-n.conn.NotifyClose(make(chan *amqp_driver.Error)))
		// Let Handle know it's not time to reconnect
		n.done <- errors.New("Channel Closed")
	}()

	log.WithField("url", n.url).Info("rbmq notifier: got connection")
	n.channel, err = n.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	n.m.Connected.Set(1)
	log.Debug("rbmq notifier got channel")
	return nil
}

func (n *Notifier) reConnect() error {

	for {
		log.WithField("reconnectDelay", n.reconnectDelay).Info("rbmq notifier reconnects...")
		time.Sleep(time.Duration(n.reconnectDelay) * time.Second)

		if err := n.connect(); err != nil {
			n.m.Connected.Set(0)
			n.m.ReconnectCount.Inc()
			log.WithField("error", err.Error()).Error("rbmq notifier could not reconnect")

		} else {
			log.WithField("url", n.url).Info("rbmq notifier connected")
			break
		}
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
	Priority  uint8
	Body      []byte
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
			if n.stop {
				return
			}

			var msg AMQPMessage
			msg, running = <-reading
			if !running {
				log.WithField("rbmq", "!running").Info("rbmq notifier")
				return
			}

			if pending <- msg; len(pending) > 0 {
				log.WithField("rbmq_pending", len(pending)).Info("rbmq notifier")
			}
		}
	}()

	for {
		if n.stop {
			return
		}
		var msg AMQPMessage
		select {
		case sig := <-n.signalsCh:
			n.stop = true
			log.WithField("signal", sig.String()).Info("rbmq notifier cought a signal")
			buf := []AMQPMessage{}
			for msg := range n.publishCh {
				buf = append(buf, msg)
			}
			for msg := range pending {
				buf = append(buf, msg)
			}
			out, err := json.Marshal(buf)
			if err != nil {
				log.WithField("pending", fmt.Sprintf("%#v", buf)).
					Error("rbmq notifier cannot marshal pending buffer")
			} else {
				fh, err := os.OpenFile(n.conf.BufferPath, os.O_CREATE|os.O_RDWR, 0744)
				if err != nil {
					log.WithField("pending", fmt.Sprintf("%#v", buf)).
						Error("rbmq notifier opern file for pending buffer")
				} else {
					fh.Write(out)
					fh.Close()
				}
			}
			n.FinishCh <- true

		case <-n.done:
			err := n.reConnect()
			if err != nil {
				// Very likely chance of failing should not cause worker to terminate
				log.WithField("error", err.Error()).Error("rbmq notifier reconnect failed")
			} else {
				log.Info("rbmq notifier: reconnected")
			}

		case msg = <-pending:
			if n.stop {
				return
			}
			q, err := n.channel.QueueDeclare(
				msg.QueueName, // name
				false,         // durable
				false,         // delete when unused
				false,         // exclusive
				false,         // no-wait
				nil,           // arguments
			)

			if err != nil {
				n.m.PublishErrs.Inc()
				n.conn.Close()

				pending <- msg
				err = fmt.Errorf("Channel.QueueDeclare: %s", err.Error())
				log.WithField("error", err.Error()).Error("rbmq notifier queue declare failed")
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
					Priority:    msg.Priority,
				})

			if err != nil {
				n.m.PublishErrs.Inc()
				n.conn.Close()

				pending <- msg
				err = fmt.Errorf("Channel.Publish: %s", err.Error())
				log.WithField("error", err.Error()).Error("rbmq notifier publish failed")
				break
			}
			log.WithFields(log.Fields{"queue": msg.QueueName, "body": string(msg.Body)}).Info("rbmq: publish")
		}
	}

}

type NotifierMetrics struct {
	SessionRequests m.Gauge
	PublishErrs     m.Gauge
	ReconnectCount  prometheus.Gauge
	Connected       prometheus.Gauge
	PendingBuffer   prometheus.Gauge
	ReadingBuffer   prometheus.Gauge
}

func newGaugeNotifier(name, help string) m.Gauge {
	return m.NewGauge("rbmq", "notifier", name, "rbmq "+help)
}
func initNotifierMetrics() NotifierMetrics {
	metrics := NotifierMetrics{
		SessionRequests: newGaugeNotifier("reconnects_count", "publisher reconnect count"),
		PublishErrs:     newGaugeNotifier("errors", "publish errors"),
		Connected:       m.PrometheusGauge("rbmq", "notifier", "connected", "publisher connection status"),
		ReconnectCount:  m.PrometheusGauge("rbmq", "notifier", "reconnect_count", "publisher connection attempts count"),
		PendingBuffer:   m.PrometheusGauge("rbmq", "notifier", "buffer_pending_gauge_size", "publisher pending buffer size"),
		ReadingBuffer:   m.PrometheusGauge("rbmq", "notifier", "buffer_reading_gauge_size", "publisher reading buffer size"),
	}
	return metrics
}
