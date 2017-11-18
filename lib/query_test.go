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

		Convey("add friends", func() {
			db.Put(
				Triple{[]byte("matteo"), []byte("friend"), []byte("daniele")},
				Triple{[]byte("daniele"), []byte("friend"), []byte("matteo")},
				Triple{[]byte("daniele"), []byte("friend"), []byte("marco")},
				Triple{[]byte("lucio"), []byte("friend"), []byte("matteo")},
				Triple{[]byte("lucio"), []byte("friend"), []byte("marco")},
				Triple{[]byte("marco"), []byte("friend"), []byte("davide")},
			)

			p := Pattern{
				Triple: Triple{
					Subject: []byte("matteo"),
				},
				SubjectV: &Variable{"x"},
			}

			r := db.VariableQuery(DefaultVariableQueryOptions, p)
			log.Println(r)
		})
	})
}
