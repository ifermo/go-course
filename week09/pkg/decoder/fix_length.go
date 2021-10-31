package decoder

import (
	"io"
	"net"
)

type FixLengthFrameDecoder struct {
	reader io.Reader
	length uint32
}

func (decoder *FixLengthFrameDecoder) Read() ([]byte, error) {
	length := int(decoder.length)
	data := make([]byte, 0, length)
	for i := 0; len(data) != length; {
		size, err := decoder.reader.Read(data[i:length])
		if err != nil {
			return nil, err
		}
		i += size
	}
	return data, nil
}

func NewFixLengthDecoder(conn net.Conn, length uint32) FrameDecoder {
	return &FixLengthFrameDecoder{
		reader: conn,
		length: length,
	}
}
