package mail

// Package mail defines the Mail struct, which encapsulates information needed to send a message.

import (
	"io"
)

// Mail represents a message to be sent, containing its URL, content type, and message body.
type Mail struct {
	Url         string    // Url is the destination URL for the message.
	ContentType string    // ContentType is the MIME type of the message body.
	Msg         io.Reader // Msg is the message body as an io.Reader.
}