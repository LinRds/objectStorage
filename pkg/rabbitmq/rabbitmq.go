package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Rabbitmq struct {
	ch *amqp.Channel
	conn *amqp.Connection
	Name string
	exchange string
}

func NewRabbitmq(rabbitServer string) *Rabbitmq {
	conn, err := amqp.Dial(rabbitServer)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	return &Rabbitmq {
		ch: ch,
		conn: conn,
		Name: q.Name,
	}
}

func (rb *Rabbitmq) Bind(exchange string) {
	err := rb.ch.QueueBind(
		rb.Name, // queue name
		"", // routing key
		exchange,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	rb.exchange = exchange
}

// Singlesend通过默认的exchange向指定的queue单发消息
func (rb *Rabbitmq) SingleSend(queue string, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	err = rb.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{ReplyTo: rb.Name, Body: b},
	)
	if err != nil {
		panic(err)
	}
}

// BroadCast向指定的exchange广播消息
func (rb *Rabbitmq) BroadCast(exchange string, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	err = rb.ch.Publish(
		exchange,
		"", // 当key为空时，该消息会群发
		false,
		false,
		amqp.Publishing{ReplyTo: rb.Name, Body: b},
	)
	if err != nil {
		panic(err)
	}
}

func (rb *Rabbitmq) Consume() (<-chan amqp.Delivery){
	ch, err := rb.ch.Consume(
		rb.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return ch
}

func (rb *Rabbitmq) Close() {
	rb.ch.Close()
	rb.conn.Close()
}