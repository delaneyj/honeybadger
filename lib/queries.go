package honeybadger

import (
	"bytes"
	"encoding/json"
	"log"
	"runtime"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

//RunQueries x
func (hb *HoneyBadger) RunQueries(queries ...QueryFn) (Solutions, error) {
	solutions := Solutions{}
	err := hb.db.View(func(tx *bolt.Tx) error {
		indexBucket := tx.Bucket(indiciesBucketName)
		if indexBucket == nil {
			return errors.New("can't load index bucket")
		}

		factsBucket := tx.Bucket(factsBucketName)
		if factsBucket == nil {
			return errors.New("can't load facts bucket")
		}

		for i, query := range queries {
			newSolutions := query(i == 0, solutions, indexBucket, factsBucket)

			solutions = newSolutions
		}
		return nil
	})

	return solutions, err
}

//QueryFn x
type QueryFn func(first bool, previousSolutions Solutions, indexBucket, factsBucket *bolt.Bucket) Solutions

func (hb *HoneyBadger) find(it indexType, variable string, aValue, bValue []byte) QueryFn {
	return func(first bool, previousSolutions Solutions, indexBucket, factsBucket *bolt.Bucket) Solutions {
		solutions := Solutions{}
		if first {
			previousSolutions = append(previousSolutions, Solution{})
		}

		for _, solution := range previousSolutions {
			var v uint64

			variableValue, ok := solution[variable]
			if ok {
				v = generateHash(variableValue)
			}

			a := generateHash(aValue)
			b := generateHash(bValue)
			prefix := keyBytes(it, a, b, v)

			c := indexBucket.Cursor()
			for index, factUUID := c.Seek(prefix); index != nil && bytes.HasPrefix(index, prefix); index, factUUID = c.Next() {
				factBytes := factsBucket.Get(factUUID)
				var f Fact
				json.Unmarshal(factBytes, &f)

				newSolution := Solution{}
				for k, v := range solution {
					newSolution[k] = v
				}

				var value []byte
				switch it {
				case spo:
					value = f.Object
				case sop:
					value = f.Predicate
				case pos:
					value = f.Subject
				default:
					log.Fatal(it, value)
				}
				newSolution[variable] = value
				solutions = append(solutions, newSolution)
			}
		}

		return solutions
	}
}

//FindSubjects x
func (hb *HoneyBadger) FindSubjects(subjectVariable string, predicate, object []byte) QueryFn {
	return hb.find(pos, subjectVariable, predicate, object)
}

//FindPredicates x
func (hb *HoneyBadger) FindPredicates(subject []byte, predicateVariable string, object []byte) QueryFn {
	return hb.find(sop, predicateVariable, subject, object)
}

//FindObjects x
func (hb *HoneyBadger) FindObjects(subject, predicate []byte, objectVariable string) QueryFn {
	return hb.find(spo, objectVariable, subject, predicate)
}

func (hb *HoneyBadger) from(it indexType, knownValue []byte, aVar, bVar string) QueryFn {
	return func(first bool, previousSolutions Solutions, indexBucket, factsBucket *bolt.Bucket) Solutions {
		solutions := Solutions{}
		if first {
			previousSolutions = append(previousSolutions, Solution{})
		}

		for _, solution := range previousSolutions {
			var a, b uint64

			aValue, ok := solution[aVar]
			if ok {
				a = generateHash(aValue)
			}

			bValue, ok := solution[bVar]
			if ok {
				b = generateHash(bValue)
			}

			k := generateHash(knownValue)
			prefix := keyBytes(it, k, a, b)

			c := indexBucket.Cursor()
			for index, factUUID := c.Seek(prefix); index != nil && bytes.HasPrefix(index, prefix); index, factUUID = c.Next() {
				factBytes := factsBucket.Get(factUUID)
				var f Fact
				json.Unmarshal(factBytes, &f)

				newSolution := Solution{}
				for k, v := range solution {
					newSolution[k] = v
				}

				switch it {
				case spo:
					newSolution[aVar] = f.Predicate
					newSolution[bVar] = f.Object
				case pso:
					newSolution[aVar] = f.Subject
					newSolution[bVar] = f.Object
				case osp:
					newSolution[aVar] = f.Subject
					newSolution[bVar] = f.Predicate
				default:
					log.Println("from index not known:", it)
					runtime.Breakpoint()
				}

				solutions = append(solutions, newSolution)
			}
		}

		return solutions
	}
}

//FromSubjects x
func (hb *HoneyBadger) FromSubjects(subject []byte, predicateVariable, objectVariable string) QueryFn {
	return hb.from(spo, subject, predicateVariable, objectVariable)
}

//FromPredicates x
func (hb *HoneyBadger) FromPredicates(subjectVariable string, predicate []byte, objectVariable string) QueryFn {
	return hb.from(pso, predicate, subjectVariable, objectVariable)
}

//FromObjects x
func (hb *HoneyBadger) FromObjects(subjectVariable, predicateVariable string, object []byte) QueryFn {
	return hb.from(osp, object, subjectVariable, predicateVariable)
}

//Solution x
type Solution map[string][]byte

//Solutions x
type Solutions []Solution
