package honeybadger

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHoneyBadger(t *testing.T) {
	Convey("With database", t, func() {
		Reset(func() {
			os.Remove("honey.db")
		})

		matteo := []byte("matteo")
		daniele := []byte("daniele")
		marco := []byte("marco")
		lucio := []byte("lucio")
		davide := []byte("davide")
		friend := []byte("friend")

		facts := Facts{
			Fact{matteo, friend, daniele},
			Fact{daniele, friend, matteo},
			Fact{daniele, friend, marco},
			Fact{lucio, friend, matteo},
			Fact{lucio, friend, marco},
			Fact{marco, friend, davide},
		}

		db, err := NewHoneyBadger()
		So(err, ShouldBeNil)
		defer db.Close()

		Convey("get all with no facts", func() {
			facts, err := db.All()
			So(facts, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "can't get facts: no facts")
		})

		Convey("add friends", func() {
			Convey("delete all facts", func() {
				err := db.Put(facts...)
				So(err, ShouldBeNil)

				err = db.Delete(facts...)
				So(err, ShouldBeNil)
				facts, err := db.All()
				So(facts, ShouldHaveLength, 0)
			})

			Convey("get all facts", func() {
				err := db.Put(facts...)
				So(err, ShouldBeNil)

				facts, err := db.All()
				So(err, ShouldBeNil)
				So(facts, ShouldHaveLength, len(facts))
			})

			x := "x"
			y := "y"

			Convey("simple queries", func() {
				db.Put(facts...)

				solutions, err := db.RunQueries(
					db.FindObjects(matteo, friend, x),
				)
				So(err, ShouldBeNil)
				So(solutions, ShouldHaveLength, 1)
				So(solutions[0], ShouldResemble, Solution{x: daniele})
			})

			Convey("joined queries", func() {
				db.Put(facts...)

				solutions, err := db.RunQueries(
					db.FindObjects(matteo, friend, x),
					db.FromPredicates(x, friend, y),
				)

				So(err, ShouldBeNil)
				So(solutions, ShouldHaveLength, 2)
				So(solutions[0], ShouldResemble, Solution{x: daniele, y: matteo})
				So(solutions[1], ShouldResemble, Solution{x: daniele, y: marco})
			})

			Convey("3 joined queries", func() {
				db.Put(facts...)

				solutions, err := db.RunQueries(
					db.FindObjects(matteo, friend, x),
					db.FromPredicates(x, friend, y),
					db.FindSubjects(y, friend, davide),
				)

				So(err, ShouldBeNil)
				So(solutions, ShouldHaveLength, 1)
				So(solutions[0], ShouldResemble, Solution{x: daniele, y: marco})
			})
		})
	})
}
