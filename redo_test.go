package redo

import (
	"testing"
	"time"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedo(t *testing.T) {
	Convey("Redo", t, func() {

		Convey(".Do", func() {
			Convey("Should successfully rerun a function that always returns error", func() {
				count := 0
				redo := NewRedo()
				err := redo.Do(3, 100*time.Millisecond, func(stop func()) error {
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
				redo := NewRedo()
				err := redo.Do(3, 100*time.Millisecond, func(stop func()) error {
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
				redo := NewRedo()
				go func() {
					err = redo.Do(3, 100*time.Millisecond, func(stop func()) error {
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

		Convey(".Backoff", func() {

			Convey("Should successfully run a function at least once and return error returned by the operation", func() {
				count := 0
				redo := NewRedo()
				bfc := NewDefaultBackoffConfig()
				bfc.MaxElapsedTime = 1 * time.Second
				err := redo.BackOff(bfc, func(stop func()) error {
					count++
					return fmt.Errorf("some error")
				})
				So(count, ShouldBeGreaterThanOrEqualTo, 1)
				So(err, ShouldNotBeNil)
			})

			Convey("Should successfully stop the execution after calling stop method passed to the function", func() {
				count := 0
				redo := NewRedo()
				bfc := NewDefaultBackoffConfig()
				bfc.MaxElapsedTime = 100 * time.Millisecond
				err := redo.BackOff(bfc, func(stop func()) error {
					count++
					stop()
					return fmt.Errorf("some error")
				})
				So(count, ShouldEqual, 1)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "some error")
			})

			Convey("Should successfully stop the execution after calling the stop method on the object", func() {
				var err error
				count := 0
				redo := NewRedo()
				bfc := NewDefaultBackoffConfig()
				bfc.MaxElapsedTime = 300 * time.Millisecond
				go func() {
					err = redo.BackOff(bfc, func(stop func()) error {
						count++
						return fmt.Errorf("some error")
					})
				}()
				time.Sleep(150 * time.Millisecond)
				redo.Stop()
				So(count, ShouldBeGreaterThanOrEqualTo, 1)
				So(err, ShouldBeNil)
				So(redo.LastErr, ShouldNotBeNil)
			})
		})

	})
}
