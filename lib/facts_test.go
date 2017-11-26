package honeybadger

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFacts(t *testing.T) {
	d := []byte("delaney")
	f := []byte("friend")
	g := []byte("guy")

	Convey("siphash values should be consistent", t, func() {
		So(generateHash(d), ShouldEqual, uint64(0xc1259675afe40a23))
		So(generateHash(f), ShouldEqual, uint64(0x70a1f82d4cb80079))
		So(generateHash(g), ShouldEqual, uint64(0xa14d59c5a214b12c))
	})

	Convey("generating index keys should give 6 keys", t, func() {
		f := Fact{d, f, g}

		keys := [][]uint8{
			{0, 35, 10, 228, 175, 117, 150, 37, 193, 121, 0, 184, 76, 45, 248, 161, 112, 44, 177, 20, 162, 197, 89, 77, 161},
			{1, 35, 10, 228, 175, 117, 150, 37, 193, 44, 177, 20, 162, 197, 89, 77, 161, 121, 0, 184, 76, 45, 248, 161, 112},
			{2, 121, 0, 184, 76, 45, 248, 161, 112, 35, 10, 228, 175, 117, 150, 37, 193, 44, 177, 20, 162, 197, 89, 77, 161},
			{3, 121, 0, 184, 76, 45, 248, 161, 112, 44, 177, 20, 162, 197, 89, 77, 161, 35, 10, 228, 175, 117, 150, 37, 193},
			{4, 44, 177, 20, 162, 197, 89, 77, 161, 35, 10, 228, 175, 117, 150, 37, 193, 121, 0, 184, 76, 45, 248, 161, 112},
			{5, 44, 177, 20, 162, 197, 89, 77, 161, 121, 0, 184, 76, 45, 248, 161, 112, 35, 10, 228, 175, 117, 150, 37, 193},
		}

		So(f.generateIndexKeys(), ShouldResemble, keys)
	})
}
