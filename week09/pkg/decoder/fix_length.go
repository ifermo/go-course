package decoder

import (
	"net"
)

type FixLengthFrameDecoder struct {
	conn   net.Conn
	length uint32
}

func (decoder *FixLengthFrameDecoder) Read() ([]byte, error) {
	length := int(decoder.length)
	data := make([]byte, 0, length)
	for i := 0; len(data) != length; {
		size, err := decoder.conn.Read(data[i:length])
		if err != nil {
			return nil, err
		}
		i += size
	}
	return data, nil
}

func NewFixLengthDecoder(conn net.Conn, length uint32) (FrameDecoder, error) {
	return &FixLengthFrameDecoder{
		conn:   conn,
		length: length,
	}, nil
}
