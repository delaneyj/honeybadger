package honeybadger

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVariable(t *testing.T) {
	SkipConvey("should test variable", t, func() {
		Convey("should have a name", func() {
			v := Variable{"x"}
			So(v.Name, ShouldEqual, "x")
		})

		Convey("should have a name (bis)", func() {
			v := Variable{"y"}
			So(v.Name, ShouldEqual, "y")
		})

		instance := Variable{"x"}

		hello := []byte("hello")
		world := []byte("world")

		Convey("isBound()", func() {
			Convey("should return true if there is a key in the solution", func() {
				So(instance.isBound(Solutions{"x": hello}), ShouldBeTrue)
			})

			Convey("should return false if there is no key in the solution", func() {
				So(instance.isBound(Solutions{}), ShouldBeFalse)
			})

			Convey("should return false if there is another key in the solution", func() {
				So(instance.isBound(Solutions{"hello": world}), ShouldBeFalse)
			})
		})

		Convey("bind()", func() {
			Convey("should return a different object", func() {
				s := Solutions{}
				So(instance.bind(s, hello), ShouldNotEqual, s)
			})

			Convey("should set an element in the solution", func() {
				s := Solutions{}
				deepEqual := Solutions{"x": hello}

				So(instance.bind(s, hello), ShouldMapContainOriginal, deepEqual)
			})

			Convey("should copy values", func() {
				s := Solutions{"y": world}
				expected := Solutions{"x": hello, "y": world}
				So(instance.bind(s, hello), ShouldMapContainOriginal, expected)
			})
		})

		Convey("isBindable()", func() {
			s := Solutions{"x": hello}

			Convey("should bind to the same value", func() {

				So(instance.isBindable(s, hello), ShouldBeTrue)
			})

			Convey("should not bind to a different value", func() {
				So(instance.isBindable(s, []byte("hello2")), ShouldBeFalse)
			})

			Convey("should bind if the key is not present", func() {
				So(instance.isBindable(Solutions{}, hello), ShouldBeTrue)
			})
		})
	})
}

func ShouldMapContainOriginal(actual interface{}, expected ...interface{}) string {
	act := actual.(Solutions)
	exp := make([]Solutions, len(expected))
	for i, x := range expected {
		exp[i] = x.(Solutions)
	}

	for _, e := range exp {
		if len(act) != len(e) {
			return fmt.Sprintf("Mismatch actual:%+v expected:%+v", actual, expected)
		}

		for k, v := range act {
			v2, ok := e[k]
			if !ok {
				return fmt.Sprintf("Missing key '%s'", k)
			}

			if !bytes.Equal(v, v2) {
				return fmt.Sprintf("Value mismatch actual:%+v expected:%+v", v, v2)
			}
		}
	}

	return ""
}
