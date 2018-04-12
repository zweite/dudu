package event

import (
	"os"
	"syscall"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEvent(t *testing.T) {
	Convey("test event", t, func() {
		Convey("wait hook", func() {
			i := 0
			hightHookFunc := func() {
				t.Log("hight hook")
				So(i, ShouldEqual, 0)
				i += 1
			}

			middleHookFunc := func() {
				t.Log("middle hook")
				So(i, ShouldEqual, 1)
				i += 2
			}

			lowHookFunc := func() {
				t.Log("low hook")
				So(i, ShouldEqual, 3)
			}

			AddHook(HightPriority, hightHookFunc)
			AddHook(MiddlePriority, middleHookFunc)
			AddHook(LowPriority, lowHookFunc)

			go func() {
				ticker := time.NewTicker(time.Second * 3)
				<-ticker.C
				if proc, err := os.FindProcess(os.Getpid()); err != nil {
					os.Exit(1)
				} else {
					proc.Signal(syscall.SIGHUP)
				}
			}()

			WaitExit()
		})
	})
}
