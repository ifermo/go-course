package internal

import (
	"encoding/binary"
	"errors"
)

const (
	PackSize      = 4
	headerSize    = 2
	verSize       = 2
	opSize        = 4
	seqSize       = 4
	rawHeaderSize = headerSize + verSize + opSize + seqSize

	// offset
	headerOffset = 0
	verOffset    = headerOffset + headerSize
	opOffset     = verOffset + verSize
	seqOffset    = opOffset + opSize
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

type Message struct {
	Version    int32
	Operation  int32
	SequenceId int32
	Data       []byte
}

func (msg *Message) Read(buf []byte) error {
	headerLen := binary.BigEndian.Uint16(buf[headerOffset:verOffset])

	msg.Version = int32(binary.BigEndian.Uint16(buf[verOffset:opOffset]))
	msg.Operation = int32(binary.BigEndian.Uint32(buf[opOffset:seqOffset]))
	msg.SequenceId = int32(binary.BigEndian.Uint32(buf[seqOffset:rawHeaderSize]))
	if headerLen != rawHeaderSize {
		return ErrProtoHeaderLen
	}
	msg.Data = buf[rawHeaderSize:]
	return nil
}

// WriteTCP write a proto to TCP writer.
func (msg *Message) Write() []byte {
	buf := make([]byte, 0, rawHeaderSize+len(msg.Data))
	binary.BigEndian.PutUint16(buf[headerOffset:], rawHeaderSize)
	binary.BigEndian.PutUint16(buf[verOffset:], uint16(msg.Version))
	binary.BigEndian.PutUint16(buf[opOffset:], uint16(msg.Operation))
	binary.BigEndian.PutUint32(buf[seqOffset:], uint32(msg.SequenceId))
	buf = append(buf, msg.Data...)
	return buf
}
