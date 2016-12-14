package amqp

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	amqp_driver "github.com/streadway/amqp"

	"github.com/vostrok/utils/metrics"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func InitQueue(
	consumer *Consumer,
	deliveryChan <-chan amqp_driver.Delivery,
	fn func(<-chan amqp_driver.Delivery),
	threads int,
	queue string,
	routingKey string,
) {
	deliveryChan, err := consumer.AnnounceQueue(queue, routingKey)
	if err != nil {
		log.WithFields(log.Fields{
			"queue": queue,
			"error": err.Error(),
		}).Fatal("rbmq consumer: AnnounceQueue")
	}
	go consumer.Handle(deliveryChan, fn, threads, queue, routingKey)
	log.WithFields(log.Fields{
		"queue": queue,
		"fn":    runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(),
	}).Info("consume init done")
}

type ConsumerMetrics struct {
	Connected          prometheus.Gauge
	ReconnectCount     prometheus.Gauge
	AnnounceQueueError prometheus.Gauge
	QueueSize          prometheus.Gauge
}

func newGaugeConsumer(name, help string) prometheus.Gauge {
	return metrics.PrometheusGauge("rbmq", "consumer", name, "rbmq consumer "+help)
}

func initConsumerMetrics(prefix string) ConsumerMetrics {
	if prefix == "" {
		log.Fatal("metrics prefix required")
	}
	return ConsumerMetrics{
		Connected:          newGaugeConsumer(prefix+"_connected", "connected"),
		ReconnectCount:     newGaugeConsumer(prefix+"_reconnect_count", "reconnect count"),
		AnnounceQueueError: newGaugeConsumer(prefix+"_announce_errors", "announce errors"),
		QueueSize:          newGaugeConsumer(prefix+"_queue_size", prefix+" queue size"),
	}
}

type ConsumerConfig struct {
	Conn           ConnectionConfig `yaml:"conn"`
	BindingKey     string           `default:"" yaml:"binding_key"`
	ExchangeType   string           `default:"" yaml:"exchange_type"`
	Exchange       string           `default:"" yaml:"exchange"`
	ReconnectDelay int              `default:"30" yaml:"reconnect_delay"`
}

type Consumer struct {
	m                  ConsumerMetrics
	queuePrefetchCount int
	conn               *amqp_driver.Connection
	channel            *amqp_driver.Channel
	done               chan error
	url                string
	exchange           string // exchange that we will bind to
	exchangeType       string // topic, direct, etc...
	bindingKey         string // routing key that we are using
	reconnectDelay     int
}

func NewConsumer(conf ConsumerConfig, queueName string, prefetchCount int) *Consumer {
	log.SetLevel(log.DebugLevel)
	url := fmt.Sprintf("amqp://%s:%s@%s:%s",
		conf.Conn.User,
		conf.Conn.Pass,
		conf.Conn.Host,
		conf.Conn.Port)

	c := &Consumer{
		m:                  initConsumerMetrics(queueName),
		queuePrefetchCount: prefetchCount,
		conn:               nil,
		channel:            nil,
		done:               make(chan error),
		url:                url,
		exchange:           conf.Exchange,
		exchangeType:       conf.ExchangeType,
		bindingKey:         conf.BindingKey,
		reconnectDelay:     conf.ReconnectDelay,
	}
	log.WithField("consumer", conf).Info("consumer init")

	go func() {
		for range time.Tick(time.Minute) {
			queueSize, err := c.GetQueueSize(queueName)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Error("cannot get queue size")
			} else {
				c.m.QueueSize.Set(float64(queueSize))
			}
		}
	}()
	return c
}

func (c *Consumer) ReConnect(queueName, bindingKey string) (<-chan amqp_driver.Delivery, error) {

	for true {
		time.Sleep(time.Duration(c.reconnectDelay) * time.Second)
		if err := c.Connect(); err != nil {
			c.m.Connected.Set(0)
			c.m.ReconnectCount.Inc()

			log.WithFields(log.Fields{
				"error":          err.Error(),
				"reconnectDelay": c.reconnectDelay,
			}).Error("consumer reconnect error")
		} else {
			log.WithFields(log.Fields{
				"url": c.url,
			}).Info("consumer connected")
			break
		}
	}

	c.m.Connected.Set(1)
	c.m.ReconnectCount.Set(0)

	deliveries, err := c.AnnounceQueue(queueName, bindingKey)
	if err != nil {
		c.m.AnnounceQueueError.Inc()

		log.WithField("error", err.Error()).Error("Could not Anounce Queue")
		return deliveries, fmt.Errorf("AnnounceQueue: %s", err.Error())
	}
	c.m.AnnounceQueueError.Set(0)
	return deliveries, nil
}

// Handle has all the logic to make sure your program keeps running
// deliveryChan should be a delievery channel as created when you call AnnounceQueue
// fn should be a function that handles the processing of deliveries
// this should be the last thing called in main as code under it will
// become unreachable unless put int a goroutine.
// The q and rk params allow you to have multiple queue listeners in main
// without them you would be tied into only using one queue per connection
func (c *Consumer) Handle(
	deliveryChan <-chan amqp_driver.Delivery,
	fn func(<-chan amqp_driver.Delivery),
	threads int,
	queue string,
	routingKey string,
) {

	var err error

	for {
		for i := 0; i < threads; i++ {
			go fn(deliveryChan)
		}

		// Go into reconnect loop when
		// c.done is passed non nil values
		if <-c.done != nil {
			deliveryChan, err = c.ReConnect(queue, routingKey)
			if err != nil {
				// Very likely chance of failing
				// should not cause worker to terminate
				log.WithField("error", err.Error()).Error("rbmq consumer reconnect failed")
			} else {
				log.Info("rbmq consumer: reconnected")
			}
		}
	}
}

// Connect to RabbitMQ server
func (c *Consumer) Connect() error {

	var err error

	log.WithField("url", c.url).Debug("dialing")
	c.conn, err = amqp_driver.Dial(c.url)
	if err != nil {
		return fmt.Errorf("amqp_driver.Dial: %s", err)
	}

	go func() {
		// Waits here for the channel to be closed
		log.Info("rbmq consumer closing: ", <-c.conn.NotifyClose(make(chan *amqp_driver.Error)))
		// Let Handle know it's not time to reconnect
		c.done <- errors.New("Channel Closed")
	}()

	log.Info("rbmq consumer: got connection, getting channel...")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	log.Info("rbmq consumer: got channel")
	//log.Info("rbmq consumer: declaring exchange (%q)", c.exchange)
	//if err = c.channel.ExchangeDeclare(
	//	c.exchange,     // name of the exchange
	//	c.exchangeType, // type
	//	true,           // durable
	//	false,          // delete when complete
	//	false,          // internal
	//	false,          // noWait
	//	nil,            // arguments
	//); err != nil {
	//	return fmt.Errorf("rbmq consumer: exchange declare: %s", err)
	//}

	c.m.Connected.Set(1)
	return nil
}

// AnnounceQueue sets the queue that will be listened to for this connection
func (c *Consumer) AnnounceQueue(queueName, bindingKey string) (<-chan amqp_driver.Delivery, error) {
	log.WithFields(log.Fields{"queue": queueName, "bindKey": bindingKey}).Debug("rbmq consumer: queue anounce")

	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		log.WithFields(log.Fields{
			"queue":   queueName,
			"bindKey": bindingKey,
			"error":   err.Error(),
		}).Error("rbmq consumer: queue declare")
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}
	log.WithFields(log.Fields{
		"queue":         queueName,
		"bindKey":       bindingKey,
		"prefetchCount": c.queuePrefetchCount,
	}).Debug("rbmq consumer: set prefetch count")

	// Qos determines the amount of messages that the queue will pass to you before
	// it waits for you to ack them. This will slow down queue consumption but
	// give you more certainty that all messages are being processed. As load increases
	// I would reccomend upping the about of Threads and Processors the go process
	// uses before changing this although you will eventually need to reach some
	// balance between threads, procs, and Qos.
	err = c.channel.Qos(c.queuePrefetchCount, 0, false)
	if err != nil {
		log.WithFields(log.Fields{
			"queue":   queueName,
			"bindKey": bindingKey,
			"error":   err.Error(),
		}).Error("rbmq consumer: set qos")
		return nil, fmt.Errorf("Error setting qos: %s", err)
	}

	//log.WithFields(log.Fields{
	//	"queue":       queue.Name,
	//	"messagesCnt": queue.Messages,
	//	"consumers":   queue.Consumers,
	//	"bindKey":     bindingKey,
	//}).Debug("rbmq soncumer: binding to exchange (key %q)")
	//
	//if err = c.channel.QueueBind(
	//	queue.Name, // name of the queue
	//	bindingKey, // bindingKey
	//	c.exchange, // sourceExchange
	//	false,      // noWait
	//	nil,        // arguments
	//); err != nil {
	//	return nil, fmt.Errorf("Queue Bind: %s", err)
	//}

	log.WithFields(log.Fields{
		"queue":         queueName,
		"bindKey":       bindingKey,
		"prefetchCount": c.queuePrefetchCount,
	}).Info("rbmq consumer: starting consume")

	deliveries, err := c.channel.Consume(
		queue.Name, // name
		"",         // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		log.WithFields(log.Fields{
			"queue":   queueName,
			"bindKey": bindingKey,
			"error":   err.Error(),
		}).Error("rbmq consumer: channel consume error")
		return nil, fmt.Errorf("rbmq consumer: queue consume: %s", err)
	}
	return deliveries, nil
}

func (c *Consumer) GetQueueSize(queue string) (int, error) {
	queueInfo, err := c.channel.QueueInspect(queue)
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
