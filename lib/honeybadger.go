package honeybadger

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	indiciesBucketName = []byte("indexes")
	factsBucketName    = []byte("facts")
)

//HoneyBadger x
type HoneyBadger struct {
	db *bolt.DB
}

//NewHoneyBadger x
func NewHoneyBadger() (*HoneyBadger, error) {
	db, err := bolt.Open("honey.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Wrap(err, "Can't start db")
	}

	hb := &HoneyBadger{db}
	return hb, nil
}

//Close x
func (hb *HoneyBadger) Close() {
	if hb.db != nil {
		hb.db.Close()
	}
}

//Put x
func (hb *HoneyBadger) Put(facts ...Fact) error {
	err := hb.db.Update(func(tx *bolt.Tx) error {
		indexBucket, err := tx.CreateBucketIfNotExists(indiciesBucketName)
		if err != nil {
			return errors.Wrap(err, "can't open or create index bucket")
		}

		factsBucket, err := tx.CreateBucketIfNotExists(factsBucketName)
		if err != nil {
			return errors.Wrap(err, "can't open or create facts bucket")
		}

		for _, f := range facts {
			//Generate all the keys and handle errors first
			factID := uuid.NewV4()
			factIDBytes := factID.Bytes()
			jsonBytes, err := json.Marshal(f)
			if err != nil {
				return errors.Wrap(err, "can't marshal the fact")
			}
			indexKeys := f.generateIndexKeys()

			//Put them in bolt
			factsBucket.Put(factIDBytes, jsonBytes)
			for _, k := range indexKeys {
				err := indexBucket.Put(k, factIDBytes)
				if err != nil {
					return errors.Wrapf(err, "can't put index in database.")
				}
			}

		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "can't put facts")
	}
	return nil
}

//Delete x
func (hb *HoneyBadger) Delete(facts ...Fact) error {
	err := hb.db.Update(func(tx *bolt.Tx) error {
		indexBucket := tx.Bucket(indiciesBucketName)
		if indexBucket == nil {
			return nil
		}

		factsBucket := tx.Bucket(factsBucketName)
		if factsBucket == nil {
			return nil
		}

		for _, f := range facts {
			indexKeys := f.generateIndexKeys()

			for i, indexKey := range indexKeys {
				if i == 0 {
					factUUID := indexBucket.Get(indexKey)
					err := factsBucket.Delete(factUUID)
					if err != nil {
						return errors.Wrap(err, "can't delete fact")
					}
				}

				err := indexBucket.Delete(indexKey)
				if err != nil {
					return errors.Wrap(err, "can't delete index")
				}
			}
		}

		return nil
	})

	if err != nil {
		err = errors.Wrap(err, "can't put facts")
	}
	return err
}

//All x
func (hb *HoneyBadger) All() (Facts, error) {
	facts := Facts{}

	err := hb.db.View(func(tx *bolt.Tx) error {
		factsBucket := tx.Bucket(factsBucketName)
		if factsBucket == nil {
			return errors.New("no facts")
		}

		c := factsBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var f Fact
			err := json.Unmarshal(v, &f)
			if err != nil {
				return errors.Wrap(err, "can't unmarshal fact")
			}
			facts = append(facts, f)
		}

		return nil
	})

	if err != nil {
		err = errors.Wrap(err, "can't get facts")
	}

	return facts, err
}
