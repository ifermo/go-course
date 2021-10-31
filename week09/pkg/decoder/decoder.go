package decoder

type FrameDecoder interface {
	Read() ([]byte, error)
}
