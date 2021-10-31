package decoder

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type ByteOrder int

const (
	BigEndian ByteOrder = iota
	LittleEndian
)

type LengthFieldBasedFrameDecoder struct {
	reader            io.Reader
	byteOrder         ByteOrder
	maxLength         uint32
	lengthFieldLength uint32
}

func (decoder *LengthFieldBasedFrameDecoder) Read() ([]byte, error) {
	frameLength, err := decoder.frameLength()
	if err != nil {
		return nil, err
	}
	data := make([]byte, 0, frameLength)
	for i := 0; len(data) < int(frameLength); {
		size, err := decoder.reader.Read(data[i:])
		if err != nil {
			return nil, err
		}
		i += size
	}
	return data, nil
}

func (decoder *LengthFieldBasedFrameDecoder) frameLength() (uint32, error) {
	lenBuf := make([]byte, 0, decoder.lengthFieldLength)
	for i := 0; len(lenBuf) < int(decoder.lengthFieldLength); {
		size, err := decoder.reader.Read(lenBuf[i:])
		if err != nil {
			return 0, err
		}
		i += size
	}
	frameLength := bytesToUint32(decoder.byteOrder, lenBuf)
	if frameLength <= decoder.lengthFieldLength || frameLength > decoder.maxLength {
		return 0, errors.New("invalid data frame length")
	}
	return frameLength, nil
}

func bytesToUint32(byteOrder ByteOrder, lenBuf []byte) uint32 {
	switch byteOrder {
	case BigEndian:
		return binary.BigEndian.Uint32(lenBuf)
	case LittleEndian:
		return binary.LittleEndian.Uint32(lenBuf)
	}
	return 0
}

func NewLengthFieldBaseFrameDecoder(conn net.Conn, byteOrder ByteOrder, maxLength uint32, lengthFieldLength uint32) FrameDecoder {
	if lengthFieldLength != 1 && lengthFieldLength != 2 &&
		lengthFieldLength != 3 && lengthFieldLength != 4 {
		panic("lengthFieldLength must be either 1, 2, 3 ,or 4")
	}
	return &LengthFieldBasedFrameDecoder{
		reader:            conn,
		byteOrder:         byteOrder,
		maxLength:         maxLength,
		lengthFieldLength: lengthFieldLength,
	}
}
