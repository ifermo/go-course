package encoder

import (
	"encoding/binary"
	"io"
	"net"
	. "week09/pkg/decoder"
)

type LengthFieldPrepender struct {
	writer            io.Writer
	byteOrder         ByteOrder
	lengthFieldLength uint32
}

func (encoder *LengthFieldPrepender) Write(data []byte) error {
	buf := make([]byte, encoder.lengthFieldLength)
	if encoder.byteOrder == BigEndian {
		binary.BigEndian.PutUint32(buf, uint32(len(data)))
	} else {
		binary.LittleEndian.PutUint32(buf, uint32(len(data)))
	}
	_, err := encoder.writer.Write(append(buf, data...))
	return err
}

func NewLengthFieldPrepender(conn net.Conn, byteOrder ByteOrder, lengthFieldLength uint32) FrameEncoder {
	if lengthFieldLength != 1 && lengthFieldLength != 2 &&
		lengthFieldLength != 3 && lengthFieldLength != 4 {
		panic("lengthFieldLength must be either 1, 2, 3 ,or 4")
	}
	return &LengthFieldPrepender{
		writer:            conn,
		byteOrder:         byteOrder,
		lengthFieldLength: lengthFieldLength,
	}
}
