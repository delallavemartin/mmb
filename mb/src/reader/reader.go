package reader

import (
	"io"
	"io/ioutil"
	"log"
)

// CustomReader, used to convert a Reader to string
type CustomReader struct {
	Reader         io.Reader
	ReaderAsString string
}

// CustomReader Methods
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
