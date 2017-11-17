package honeybadger

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVariable(t *testing.T) {
	Convey("should test variable", t, func() {
		Convey("should have a name", func() {
			v := Variable{"x"}
			So(v.Name, ShouldEqual, "x")
		})

		Convey("should have a name (bis)", func() {
			v := Variable{"y"}
			So(v.Name, ShouldEqual, "y")
		})

		instance := Variable{"x"}

		Convey("isBound()", func() {
			Convey("should return true if there is a key in the solution", func() {
				So(instance.isBound(Solution{"x": "hello"}), ShouldBeTrue)
			})

			Convey("should return false if there is no key in the solution", func() {
				So(instance.isBound(Solution{}), ShouldBeFalse)
			})

			Convey("should return false if there is another key in the solution", func() {
				So(instance.isBound(Solution{"hello": "world"}), ShouldBeFalse)
			})
		})

		Convey("bind()", func() {
			Convey("should return a different object", func() {
				solution := Solution{}
				So(instance.bind(solution, "hello"), ShouldNotEqual, solution)
			})

			Convey("should set an element in the solution", func() {
				solution := Solution{}
				deepEqual := Solution{"x": "hello"}

				So(instance.bind(solution, "hello"), ShouldMapContainOriginal, deepEqual)
			})

			Convey("should copy values", func() {
				solution := Solution{"y": "world"}
				expected := Solution{"x": "hello", "y": "world"}
				So(instance.bind(solution, "hello"), ShouldMapContainOriginal, expected)
			})

		})

		Convey("isBindable()", func() {
			s := Solution{"x": "hello"}

			Convey("should bind to the same value", func() {

				So(instance.isBindable(s, "hello"), ShouldBeTrue)
			})

			Convey("should not bind to a different value", func() {
				So(instance.isBindable(s, "hello2"), ShouldBeFalse)
			})

			Convey("should bind if the key is not present", func() {
				So(instance.isBindable(Solution{}, "hello"), ShouldBeTrue)
			})
		})
	})
}

func ShouldMapContainOriginal(actual interface{}, expected ...interface{}) string {
	act := actual.(Solution)
	exp := make([]Solution, len(expected))
	for i, x := range expected {
		exp[i] = x.(Solution)
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

			if v != v2 {
				return fmt.Sprintf("Value mismatch actual:%+v expected:%+v", v, v2)
			}
		}
	}

	return ""
}
