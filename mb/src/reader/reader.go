package reader

// Package reader provides a custom reader utility for converting an io.Reader to a string.

import (
	"io"
	"io/ioutil"
	"log"
)

// CustomReader is a utility struct that wraps an io.Reader and provides a method
// to convert its content to a string, caching the result.
type CustomReader struct {
	Reader         io.Reader // Reader is the underlying io.Reader.
	ReaderAsString string    // ReaderAsString caches the string representation of the reader's content.
}

// ToString reads the content of the underlying Reader and returns it as a string.
// The result is cached for subsequent calls. If an error occurs during reading,
// it logs the error and returns an empty string.
func (self *CustomReader) ToString() string {
	if self.ReaderAsString == "" {
		stream, err := ioutil.ReadAll(self.Reader)
		if err != nil {
			log.Println("ERROR - reading reader FAILED")
			return ""
		}
		self.ReaderAsString = string(stream)
	}
	return self.ReaderAsString
}
