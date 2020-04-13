package messaging

import (
	"bytes"
	"encoding/gob"

	"github.com/streadway/amqp"
)

// Writer the publisher's writer
type Writer struct {
	Buf *bytes.Buffer
	Enc *gob.Encoder
}

func newWriter() *Writer {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	return &Writer{
		Buf: buf,
		Enc: enc,
	}
}

// Reset the writer
func (w *Writer) Reset() {
	w.Buf.Reset()
	enc := gob.NewEncoder(w.Buf)
	w.Enc = enc
}

// MessageBytes
func (p *Publisher) MessageBytes() []byte {
	return p.Writer.Buf.Bytes()
}

// Publisher publishes to the queue
type Publisher struct {
	Server       *Server
	Channel      *amqp.Channel
	Queue        amqp.Queue
	QueueName    string
	ExchangeName string
	Writer       *Writer
}

// NewPublisherWithQueue returns a new publisher with a declared queue
func NewPublisherWithQueue(s *Server, queue string, autoDelete bool) *Publisher {
	ch, err := s.Conn.Channel()
	FailOnError(err, "Failed to open channel")
	q, err := ch.QueueDeclare(queue, false, autoDelete, false, false, nil)
	FailOnError(err, "Failed to declare and connect to queue")
	return &Publisher{
		Server:       s,
		Queue:        q,
		QueueName:    q.Name,
		ExchangeName: "",
		Channel:      ch,
	}
}

// NewPublisherWithExchange returns a publisher struct with the server and queue
func NewPublisherWithExchange(s *Server, exchange string) *Publisher {
	ch, err := s.Conn.Channel()
	FailOnError(err, "Failed to open channel")
	return &Publisher{
		Server:       s,
		Channel:      ch,
		QueueName:    "",
		ExchangeName: exchange,
	}
}

// Message returns a message to be published
func (p *Publisher) Message(contentType string, body []byte) amqp.Publishing {
	return amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	}
}

// Publish publishes a message
func (p *Publisher) Publish(bytes []byte) {
	msg := p.Message("text/plain", bytes)
	p.Channel.Publish(p.ExchangeName, p.QueueName, false, false, msg)
}

// Stop closes the channel connection
func (p *Publisher) Stop() {
	p.Channel.Close()
}

// SetUpWriter sets up the buffer and encoder
func (p *Publisher) SetUpWriter() {
	p.Writer = newWriter()
}

// WriteMessageToBuffer  writes te message to the buffer
func (p *Publisher) WriteMessageToBuffer(message MessagInterface) {
	message.Encode(p.Writer)
}
