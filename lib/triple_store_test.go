package honeybadger

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTripleStore(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	SkipConvey("with a new DB", t, func() {
		db, err := New("")
		So(err, ShouldBeNil)

		Reset(func() {
			db.Close()
		})

		a := []byte("a")
		b := []byte("b")
		c := []byte("c")
		d := []byte("d")
		qp := QueryPattern{}

		Convey("with a triple inserted", func() {
			triple := Triple{
				Subject:   a,
				Predicate: b,
				Object:    c,
			}

			db.Put(triple)

			Convey("should get it specifiying the subject", func() {
				qp.Triple = Triple{Subject: a}
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should get it specifiying the object", func() {
				qp.Triple = Triple{Object: c}
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should get it specifiying the predicate", func() {
				qp.Triple = Triple{Predicate: b}
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should get it specifiying the subject and the predicate", func() {
				qp.Triple = Triple{Subject: a, Predicate: b}
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should get it specifiying the subject and the object", func() {
				qp.Triple = Triple{Subject: a, Object: c}
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should get it specifiying the predicate and the object", func() {
				qp.Triple = Triple{Predicate: b, Object: c}
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should return the triple through the getStream interface", func() {
				qp.Triple = Triple{Predicate: b}
				for t := range db.SearchCh(qp) {
					So(t, ShouldResemble, triple)
				}
			})

			Convey("should get the triple if limit 1 is used", func() {
				qp.Limit = 1
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should get the triple if limit 0 is used", func() {
				qp.Limit = 0
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should get the triple if offset 0 is used", func() {
				qp.Offset = 0
				So(db.Search(qp), ShouldContain, triple)
			})

			Convey("should not get the triple if offset 1 is used", func() {
				qp.Offset = 1
				So(db.Search(qp), ShouldNotContain, triple)
			})
		})

		Convey("should put an array of triples", func() {
			t1 := Triple{a, b, c}
			t2 := Triple{a, b, d}
			db.Put(t1, t2)
		})

		Convey("should get only triples with exact match of subjects", func() {
			t1 := Triple{[]byte("a1"), b, c}
			t2 := Triple{a, b, d}
			db.Put(t1, t2)

			qp.Triple = Triple{Subject: a}
			results := db.Search(qp)
			So(results, ShouldHaveLength, 1)
			So(results[0], ShouldResemble, t2)
		})

		Convey("with special characters", func() {
			Convey("should support string contain ::", func() {
				t1 := Triple{Subject: a, Predicate: b, Object: c}
				t2 := Triple{Subject: []byte("a::a::a"), Predicate: b, Object: c}
				db.Put(t1, t2)

				qp.Triple = Triple{Subject: a}
				So(db.Search(qp), ShouldHaveLength, 1)
			})

			Convey("should support string contain \\::", func() {
				t1 := Triple{Subject: a, Predicate: b, Object: c}
				t2 := Triple{Subject: []byte("a\\::a"), Predicate: b, Object: c}
				db.Put(t1, t2)

				qp.Triple = Triple{Subject: a}
				So(db.Search(qp), ShouldHaveLength, 1)
			})
			Convey("should support string end with :", func() {
				aColon := []byte("a:")
				t1 := Triple{Subject: a, Predicate: b, Object: c}
				t2 := Triple{Subject: aColon, Predicate: b, Object: c}
				db.Put(t1, t2)

				qp.Triple = Triple{Subject: aColon}
				result := db.Search(qp)
				So(result, ShouldHaveLength, 1)
				So(result[0].Subject, ShouldResemble, aColon)
			})

			Convey("should support string end with \\", func() {
				aBackslash := []byte("a\\")
				t1 := Triple{Subject: a, Predicate: b, Object: c}
				t2 := Triple{Subject: aBackslash, Predicate: b, Object: c}
				db.Put(t1, t2)

				qp.Triple = Triple{Subject: aBackslash}
				res := db.Search(qp)
				So(res, ShouldHaveLength, 1)
				So(res[0].Subject, ShouldResemble, aBackslash)
			})
		})

		Convey("should put a triple with an object to false", func() {
			f := []byte{0}
			t := Triple{Subject: a, Predicate: b, Object: f}
			db.Put(t)

			qp.Triple = t
			results := db.Search(qp)
			log.Printf("%+v", results)
			So(results, ShouldHaveLength, 1)
		})
	})
}
