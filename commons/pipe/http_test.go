package pipe

import (
	"net"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpPipe(t *testing.T) {
	tbody := []byte("hello world")
	addr := "127.0.0.1:2002"
	pattern := "/abc"
	auth := "fuck"

	listen, err := net.Listen("tcp", addr)
	require.Nil(t, err)

	pop := NewHttpPipePop(auth)
	go func() {
		// new http serv
		http.Handle(pattern, pop)
		http.Serve(listen, nil)
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		dataChan, err := pop.Pop()
		require.Nil(t, err)
		select {
		case data := <-dataChan:
			require.Equal(t, tbody, data)
			return
		}
	}()

	push := NewHttpPipePush("http://"+addr+pattern, auth)
	var pipe Pipe
	pipe = NewHttpPipe(push, pop)
	pipe.Push(tbody)
	wg.Wait()

	// 关闭资源
	pipe.Stop()
	listen.Close()
}
