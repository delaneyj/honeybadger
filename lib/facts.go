package honeybadger

import (
	"bytes"
	"encoding/binary"

	"github.com/dchest/siphash"
)

const (
	key0 = uint64(0x10d56607b894f605) // Siphash of "honeybadger" with self key
	key1 = uint64(0xFFFFFFFFFFFFFFC5) // Largest 64bit prime
)

type indexType byte

const (
	spo indexType = iota
	sop
	pso
	pos
	osp
	ops
)

//Fact  x
type Fact struct {
	Subject   []byte
	Predicate []byte
	Object    []byte
	// At        time.Time
}

//Facts x
type Facts []Fact

func (f *Fact) generateIndexKeys() [][]byte {
	s := generateHash(f.Subject)
	p := generateHash(f.Predicate)
	o := generateHash(f.Object)

	return [][]byte{
		keyBytes(spo, s, p, o),
		keyBytes(sop, s, o, p),
		keyBytes(pso, p, s, o),
		keyBytes(pos, p, o, s),
		keyBytes(osp, o, s, p),
		keyBytes(ops, o, p, s),
	}
}

func keyBytes(prefix indexType, a, b, c uint64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, prefix)

	if a > 0 {
		binary.Write(buf, binary.LittleEndian, a)
		if b > 0 {
			binary.Write(buf, binary.LittleEndian, b)
			if c > 0 {
				binary.Write(buf, binary.LittleEndian, c)
			}
		}
	}

	return buf.Bytes()
}

func generateHash(x []byte) uint64 {
	return siphash.Hash(key0, key1, x)
}
