package encoder

type FrameEncoder interface {
	Write([]byte) error
}
