package honeybadger

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	triplesPrefix = []byte("triples_")
)

//Triple x
type Triple struct {
	Subject   []byte
	Predicate []byte
	Object    []byte
}

//HoneyBadger x
type HoneyBadger struct {
	db *badger.DB
}

//New x
func New(path string) (*HoneyBadger, error) {
	opts := badger.DefaultOptions

	var dir string
	var err error
	if len(path) == 0 {
		dir, err = ioutil.TempDir("", "example")
		if err != nil {
			return nil, errors.Wrap(err, "Can't create temp dir")
		}
	}

	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	return &HoneyBadger{
		db: db,
	}, nil
}

//Close x
func (hb *HoneyBadger) Close() {
	hb.db.Close()
}

//Put x
func (hb *HoneyBadger) Put(triples ...Triple) {
	err := hb.db.Update(func(txn *badger.Txn) error {
		for _, t := range triples {
			err := put(txn, t)
			if err != nil {
				return errors.Wrap(err, "can't put triple")
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

//PutCh x
func (hb *HoneyBadger) PutCh(done chan bool) chan Triple {
	ch := make(chan Triple)

	go func() {
		err := hb.db.Update(func(txn *badger.Txn) error {
			for t := range ch {
				err := put(txn, t)
				if err != nil {
					return errors.Wrap(err, "can't put triple")
				}
			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		done <- true
	}()

	return ch
}

func put(txn *badger.Txn, t Triple) error {
	triplesID := append(triplesPrefix, uuid.NewV1().Bytes()...)
	tripleBytes, err := json.Marshal(t)
	if err != nil {
		return errors.Wrap(err, "can't marshal triple")
	}
	txn.Set(triplesID, tripleBytes)

	keys := genKeys(t)
	for _, k := range keys {
		txn.Set(k, triplesID)
	}

	return nil
}

//Delete x
func (hb *HoneyBadger) Delete(triples ...Triple) {
	err := hb.db.Update(func(tx *badger.Txn) error {
		for _, t := range triples {
			delete(tx, t)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

//DeleteCh x
func (hb *HoneyBadger) DeleteCh(done chan bool) chan Triple {
	ch := make(chan Triple)

	go func() {
		err := hb.db.Update(func(txn *badger.Txn) error {
			for t := range ch {
				err := delete(txn, t)
				if err != nil {
					return errors.Wrap(err, "can't delete triple")
				}
			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		done <- true
	}()

	return ch
}

func delete(tx *badger.Txn, t Triple) error {
	keys := genKeys(t)
	for i, key := range keys {
		if i == 0 {
			item, err := tx.Get(key)
			if err != nil {
				return errors.Wrap(err, "can't get triple item from badger")
			}
			tripleID, err := item.Value()
			if err != nil {
				return errors.Wrap(err, "can't get triple uuid")
			}
			tx.Delete(tripleID)
		}
		tx.Delete(key)
	}
	return nil
}

//Count x
func (hb *HoneyBadger) Count(pattern Pattern) (uint, uint) {
	query := createQuery(pattern, streamOptions{})
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false

	var offset, counted, rdfCount, indexBytes uint
	indexTypeCount := len(defs)
	err := hb.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(opts)
		for it.Seek(query.Prefix); it.ValidForPrefix(query.Prefix); it.Next() {
			if offset < pattern.Offset {
				offset++
				continue
			}

			if pattern.Limit != 0 && counted >= pattern.Limit {
				break
			}

			key := it.Item().Key()
			indexBytes += uint(indexTypeCount * len(key))
			rdfCount++
			counted++
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return rdfCount, indexBytes
}

//Pattern x
type Pattern struct {
	Limit   uint
	Offset  uint
	Triple  Triple
	Reverse bool
	Filter  func(t Triple) bool
}

type streamOptions struct {
	Index string
}

//Search x
func (hb *HoneyBadger) Search(pattern Pattern) []Triple {
	triples := []Triple{}
	for t := range hb.SearchCh(pattern) {
		triples = append(triples, t)
	}
	return triples
}

//SearchCh x
func (hb *HoneyBadger) SearchCh(pattern Pattern) chan Triple {
	tripleseCh := make(chan Triple)

	go func() {
		query := createQuery(pattern, streamOptions{})
		startingPrefix := query.Prefix

		err := hb.db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.Reverse = pattern.Reverse

			if opts.Reverse {
				startingPrefix = append(query.Prefix, 0xff)
			}

			var offset, sent uint
			it := txn.NewIterator(opts)

			for it.Seek(startingPrefix); it.ValidForPrefix(query.Prefix); it.Next() {
				if offset < pattern.Offset {
					offset++
					continue
				}

				if pattern.Limit != 0 && sent >= pattern.Limit {
					break
				}

				item := it.Item()
				tripleID, err := item.Value()
				if err != nil {
					return errors.Wrap(err, "can't get triple id")
				}

				tItem, err := txn.Get(tripleID)
				if err != nil {
					return errors.Wrap(err, "can't get triple contents")
				}

				bytes, err := tItem.Value()
				var triple Triple
				err = json.Unmarshal(bytes, &triple)
				if err != nil {
					return errors.Wrap(err, "can't unmarshal triple")
				}

				if pattern.Filter != nil {
					shouldKeep := pattern.Filter(triple)
					if !shouldKeep {
						continue
					}
				}

				tripleseCh <- triple
				sent++
			}
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
		close(tripleseCh)
	}()

	return tripleseCh
}

//JoinAlgorithm x
type JoinAlgorithm int

const (
	//Sort x
	Sort JoinAlgorithm = iota
)

//VariableQueryOptions x
type VariableQueryOptions struct {
	Algorithm JoinAlgorithm
}

var (
	//DefaultVariableQueryOptions x
	DefaultVariableQueryOptions = VariableQueryOptions{
		Algorithm: Sort,
	}
)

//VariableQuery x
func (hb *HoneyBadger) VariableQuery(query Query, options VariableQueryOptions) {
	planner := NewQueryPlanner(hb, options)
	// result := NewPassthrough()
	// result.ObjectMode = true
	log.Fatal(planner)
}
