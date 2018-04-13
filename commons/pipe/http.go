package pipe

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpPipe struct {
	PipePusher
	PipePoper
}

func NewHttpPipe(push *HttpPipePush, pop *HttpPipePop) *HttpPipe {
	return &HttpPipe{
		PipePusher: push,
		PipePoper:  pop,
	}
}

type HttpPipePush struct {
	auth   string // 认证信息
	addr   string // 传输地址
	client *http.Client
}

// 只支持推送
func NewHttpPipePush(addr, auth string) *HttpPipePush {
	return &HttpPipePush{
		auth:   auth,
		addr:   addr,
		client: new(http.Client),
	}
}

func (h *HttpPipePush) Push(data []byte) (err error) {
	req, err := http.NewRequest("POST", h.addr, bytes.NewReader(data))
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Basic "+h.auth)
	resp, err := h.client.Do(req)
	if err != nil {
		return
	}

	resp.Body.Close()
	return
}

type HttpPipePop struct {
	isStop chan struct{}
	auth   string
	data   chan []byte
}

func NewHttpPipePop(auth string) *HttpPipePop {
	return &HttpPipePop{
		isStop: make(chan struct{}),
		auth:   "Basic " + auth,
		data:   make(chan []byte, 10), // want buf ?
	}
}

func (h *HttpPipePop) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	select {
	case <-h.isStop:
		w.WriteHeader(http.StatusForbidden)
		return
	default:
	}

	auth := r.Header.Get("Authorization")
	if auth != h.auth {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.data <- data
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *HttpPipePop) Pop() (<-chan []byte, error) {
	return h.data, nil
}

func (h *HttpPipePop) Stop() {
	close(h.isStop)
	for {
		if len(h.data) == 0 {
			break
		}
		time.Sleep(time.Second * 3) // wait for proc most rest req
	}
	close(h.data)
}
