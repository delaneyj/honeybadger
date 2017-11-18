package honeybadger

import (
	"bytes"
	"log"
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//TestTripleStoreUnicode x
func TestTripleStoreUnicode(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	t.Parallel()

	SkipConvey("a basic unicode triple store", t, func() {
		db, err := New("")
		So(err, ShouldBeNil)

		Reset(func() {
			db.Close()
		})

		p := Pattern{}

		kanjiCar := []byte("车")
		kanjiYes := []byte("是")
		kanjiTransportation := []byte("交通工具")
		kanjiAnimal := []byte("动物")
		kanjiAircarft := []byte("飞机")
		kanjiBear := []byte("狗熊")

		Convey("should put a triple", func() {
			db.Put(Triple{kanjiCar, kanjiYes, kanjiTransportation})
		})

		Convey("with a triple inserted", func() {
			square := []byte("􀃿")
			alchemicalAirSymbol := []byte("🜁")
			railwayEmoji := []byte("🚃")
			triple := Triple{square, alchemicalAirSymbol, railwayEmoji}
			db.Put(triple)

			Convey("should get it specifiying the subject", func() {
				p.Triple.Subject = square
				So(db.Search(p), ShouldContain, triple)
			})

			Convey("should get it specifiying the object", func() {
				p.Triple.Object = railwayEmoji
				So(db.Search(p), ShouldContain, triple)
			})

			Convey("should get it specifiying the predicate", func() {
				p.Triple.Predicate = alchemicalAirSymbol
				So(db.Search(p), ShouldContain, triple)
			})

			Convey("should get it specifiying the subject and the predicate", func() {
				p.Triple.Subject = square
				p.Triple.Predicate = alchemicalAirSymbol
				So(db.Search(p), ShouldContain, triple)
			})

			Convey("should get it specifiying the subject and the object", func() {
				p.Triple.Subject = square
				p.Triple.Object = railwayEmoji
				So(db.Search(p), ShouldContain, triple)
			})

			Convey("should get it specifiying the predicate and the object", func() {
				p.Triple.Predicate = alchemicalAirSymbol
				p.Triple.Object = railwayEmoji
				So(db.Search(p), ShouldContain, triple)
			})

			Convey("should return the triple through the SearchCh interface", func() {
				p.Triple.Predicate = alchemicalAirSymbol
				result := <-db.SearchCh(p)
				So(result, ShouldResemble, triple)
			})

			Convey("should get the triple if limit 1 is used", func() {
				p.Limit = 1
				So(db.Search(p), ShouldContain, triple)
			})

			Convey("should get the triple if limit 0 is used", func() {
				p.Limit = 0
				So(db.Search(p), ShouldContain, triple)
			})

			Convey("should get the triple if offset 0 is used", func() {
				p.Offset = 0
				r := db.Search(p)
				So(r, ShouldHaveLength, 1)
				So(r, ShouldContain, triple)
			})

			Convey("should not get the triple if offset 1 is used", func() {
				p.Offset = 1
				r := db.Search(p)
				So(r, ShouldHaveLength, 0)
				So(r, ShouldNotContain, triple)
			})
		})

		Convey("should put an array of triples", func() {
			t1 := Triple{kanjiCar, kanjiYes, kanjiTransportation}
			t2 := Triple{kanjiCar, kanjiYes, kanjiAnimal}

			db.Put(t1, t2)

			Convey("should get only triples with exact match of subjects", func() {
				db.Put(Triple{kanjiAircarft, kanjiYes, kanjiTransportation})

				p.Triple.Subject = kanjiCar
				results := db.Search(p)
				So(results, ShouldHaveLength, 2)
				So(results, ShouldContain, t1)
				So(results, ShouldContain, t2)
			})
		})

		Convey("with two triple inserted with the same predicate", func() {
			t1 := Triple{kanjiAircarft, kanjiYes, kanjiTransportation}
			t2 := Triple{kanjiBear, kanjiYes, kanjiAnimal}
			db.Put(t1, t2)

			Convey("should get one by specifiying the subject", func() {
				p.Triple.Subject = kanjiAircarft
				r := db.Search(p)
				So(r, ShouldHaveLength, 1)
				So(r, ShouldContain, t1)
			})

			p.Triple.Predicate = kanjiYes
			Convey("should get two by specifiying the predicate", func() {
				r := db.Search(p)
				So(r, ShouldHaveLength, 2)
				So(r, ShouldContain, t1)
				So(r, ShouldContain, t2)
			})

			Convey("should remove one and still return the other", func() {
				db.Delete(t2)
				r := db.Search(p)
				So(r, ShouldHaveLength, 1)
				So(r, ShouldContain, t1)
				So(r, ShouldNotContain, t2)
			})

			Convey("should return both triples through the getStream interface", func() {
				ch := db.SearchCh(p)
				So(t1, ShouldResemble, <-ch)
				So(t2, ShouldResemble, <-ch)
			})

			Convey("should return only one triple with limit 1", func() {
				p.Limit = 1
				r := db.Search(p)
				So(r, ShouldHaveLength, 1)
				So(r, ShouldContain, t1)
			})

			Convey("should return two triples with limit 2", func() {
				p.Limit = 2
				r := db.Search(p)
				So(r, ShouldHaveLength, 2)
				So(r, ShouldContain, t1)
				So(r, ShouldContain, t2)
			})

			Convey("should return 2 triples with limit 3", func() {
				p.Limit = 3
				r := db.Search(p)
				So(r, ShouldHaveLength, 2)
				So(r, ShouldContain, t1)
				So(r, ShouldContain, t2)
			})

			Convey("should support limit over streams", func() {
				p.Limit = 1
				So(t1, ShouldResemble, <-db.SearchCh(p))
			})

			Convey("should return only one triple with offset 1", func() {
				p.Offset = 1
				r := db.Search(p)
				So(r, ShouldHaveLength, 1)
				So(r, ShouldContain, t2)
			})

			Convey("should return only no triples with offset 2", func() {
				p.Offset = 2
				r := db.Search(p)
				So(r, ShouldBeEmpty)
			})

			Convey("should support offset over streams", func() {
				p.Offset = 1
				So(t2, ShouldResemble, <-db.SearchCh(p))
			})

			Convey("should return the triples in reverse order with reverse true", func() {
				p.Reverse = true
				r := db.Search(p)
				So(r, ShouldHaveLength, 2)
				So(r[0], ShouldResemble, t2)
				So(r[1], ShouldResemble, t1)
			})

			Convey("should return the last triple with reverse true and limit 1", func() {
				p.Reverse = true
				p.Limit = 1
				r := db.Search(p)
				So(r, ShouldHaveLength, 1)
				So(r[0], ShouldResemble, t2)
			})

			Convey("should support reverse over streams", func() {
				p.Reverse = true
				ch := db.SearchCh(p)
				So(<-ch, ShouldResemble, t2)
				So(<-ch, ShouldResemble, t1)
			})
		})

		Convey("with 10 triples inserted", func() {
			kanjiTestItems := []byte("测试项")
			kanjiDerived := []byte("推导")
			kanjiAims := []byte("目标")

			count := 10
			triples := make([]Triple, count)
			for i := range triples {
				triples[i].Subject = kanjiTestItems
				triples[i].Predicate = kanjiDerived

				l := len(kanjiAims)
				o := make([]byte, l+1)
				copy(o[:], kanjiAims)
				o[l] = byte(i)
				triples[i].Object = o
			}
			db.Put(triples...)

			Convey("should return the approximate size", func() {
				p.Triple.Predicate = kanjiDerived
				rdfCount, indexBytes := db.Count(p)
				So(rdfCount, ShouldEqual, count)
				So(indexBytes, ShouldAlmostEqual, 3060)
			})
		})

		Convey("should put triples using a stream", func() {
			done := make(chan bool)
			t1 := Triple{kanjiCar, kanjiYes, kanjiTransportation}
			t2 := Triple{kanjiCar, kanjiYes, kanjiAnimal}
			ch := db.PutCh(done)
			ch <- t1
			ch <- t2
			close(ch)
			<-done

			p.Triple.Predicate = kanjiYes

			Convey("should store the triples written using a stream", func() {
				r := db.Search(p)
				So(r, ShouldHaveLength, 2)
				So(r[0], ShouldResemble, t1)
				So(r[1], ShouldResemble, t2)
			})

			Convey("should support filtering", func() {
				p.Filter = func(t Triple) bool {
					return bytes.Compare(t.Object, kanjiTransportation) == 0
				}
				r := db.Search(p)
				So(r, ShouldHaveLength, 1)
				So(r, ShouldContain, t1)
			})

			Convey("should del the triples using a stream", func() {
				done := make(chan bool)
				ch := db.DeleteCh(done)
				ch <- t1
				ch <- t2
				close(ch)
				<-done

				r := db.Search(p)
				So(r, ShouldBeEmpty)
			})
		})

		Convey("generateBatch", func() {
			Convey("should generate a batch from a triple", func() {
				triple := Triple{kanjiCar, kanjiYes, kanjiTransportation}
				keys := genKeys(triple)

				So(keys, ShouldHaveLength, 6)

				defKeys := make([][]byte, len(defs))
				i := 0
				for k := range defs {
					defKeys[i] = []byte(k)
					i++
				}

				sort.Slice(keys, func(i, j int) bool {
					return bytes.Compare(keys[i], keys[j]) < 0
				})

				sort.Slice(defKeys, func(i, j int) bool {
					return bytes.Compare(defKeys[i], defKeys[j]) < 0
				})

				for i, k := range keys {
					So(k[0:3], ShouldResemble, defKeys[i])
				}
			})
		})
	})
}
