package pipe

// 管道接口
// 日志传输管道
type Pipe interface {
	PipePusher
	PipePoper
}

type PipePusher interface {
	Push([]byte) error
}

type PipePoper interface {
	Pop() (<-chan []byte, error)
	Stop()
}
