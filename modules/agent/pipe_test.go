package agent

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHttpPipe(t *testing.T) {
	tbody := []byte("hello world")
	addr := "127.0.0.1:2002"
	pattern := "/abc"
	auth := "fuck"

	go func() {
		// new http serv
		http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Basic "+auth {
				t.Fail()
			}

			body, err := ioutil.ReadAll(r.Body)
			r.Body.Close()

			require.Nil(t, err)
			require.Equal(t, tbody, body)
		})

		http.ListenAndServe(addr, nil)
	}()

	var pipe Pipe
	pipe = NewHttpPipe("http://"+addr+pattern, auth)
	pipe.Push(tbody)
	time.Sleep(time.Second * 3)
}
