package mail

import (
	"io"
)

// Store all information needed to publish the request.
type Mail struct {
	Url         string
	ContentType string
	Msg         io.Reader
}