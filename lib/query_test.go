package honeybadger

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQueryPlanner(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Convey("with a new DB", t, func() {
		db, err := New("")
		So(err, ShouldBeNil)

		Reset(func() {
			db.Close()
		})

		matteo := []byte("matteo")
		daniele := []byte("daniele")
		lucio := []byte("lucio")
		marco := []byte("marco")
		davide := []byte("davide")
		friend := []byte("friend")

		Convey("add friends", func() {
			db.Put(
				Triple{matteo, friend, daniele},
				Triple{daniele, friend, matteo},
				Triple{daniele, friend, marco},
				Triple{lucio, friend, matteo},
				Triple{lucio, friend, marco},
				Triple{marco, friend, davide},
			)

			qp := QueryPattern{
				Triple: Triple{
					Subject: []byte("matteo"),
				},
				Variables: TripleVariables{
					Subject: &Variable{"x"},
				},
			}

			r, err := db.VariableQuery(DefaultVariableQueryOptions, qp)
			So(err, ShouldBeNil)

			So(r, ShouldHaveLength, 1)

			f, ok := r["x"]
			So(ok, ShouldBeTrue)
			So(f, ShouldEqual, daniele)
		})
	})
}
