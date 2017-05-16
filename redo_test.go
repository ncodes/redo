package redo

import (
	"testing"
	"time"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedo(t *testing.T) {
	Convey("Redo", t, func() {

		Convey("Should successfully rerun a function that always returns error", func() {
			count := 0
			redo := NewRedo(3, 100*time.Millisecond)
			err := redo.Do(func(stop func()) error {
				count++
				return fmt.Errorf("count %d", count)
			})
			So(count, ShouldEqual, 3)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrMaxRetryReached)
			So(redo.LastErr.Error(), ShouldEqual, "count 3")
		})

		Convey("Should successfully stop the rerun of the function by calling the stop function passed to the running function", func() {
			count := 0
			redo := NewRedo(3, 100*time.Millisecond)
			err := redo.Do(func(stop func()) error {
				if count == 2 {
					stop()
					return fmt.Errorf("stop called from func")
				}
				count++
				return fmt.Errorf("something bad")
			})
			So(count, ShouldEqual, 2)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "stop called from func")
		})

		Convey("Should successfully stop the rerun of the function by calling the stop method of the object", func() {
			var err error
			count := 0
			redo := NewRedo(3, 100*time.Millisecond)
			go func() {
				err = redo.Do(func(stop func()) error {
					time.Sleep(100 * time.Millisecond)
					count++
					return fmt.Errorf("something bad")
				})
			}()
			time.Sleep(150 * time.Millisecond)
			redo.Stop()
			So(count, ShouldEqual, 1)
		})

	})
}
