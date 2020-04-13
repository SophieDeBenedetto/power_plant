package messaging

// MessagInterface the message interface
type MessagInterface interface {
	Encode(writer *Writer)
}
