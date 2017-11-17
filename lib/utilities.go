package honeybadger

import (
	"hash"
	"log"
	"runtime"
	"sort"

	"github.com/aead/siphash"
)

var (
	defs = map[string][]string{
		"spo": []string{"subject", "predicate", "object"},
		"sop": []string{"subject", "object", "predicate"},
		"pos": []string{"predicate", "object", "subject"},
		"pso": []string{"predicate", "subject", "object"},
		"ops": []string{"object", "predicate", "subject"},
		"osp": []string{"object", "subject", "predicate"},
	}
	h hash.Hash
)

type criteriaFn func(triple Triple, key string) bool

func objectMask(criteria criteriaFn, object Query) Pattern {
	p := Pattern{}
	log.Fatal()
	// 	valid := false
	// loop:
	// 	for k, v := range object {
	// 		for _, x := range defs["spo"] {
	// 			if k == x {
	// 				break loop
	// 			}
	// 		}

	// 		if criteria(object, k) {
	// 			break
	// 		}

	// 	}
	// 	//     return lodashKeys(object).
	// 	// 	filter(a).
	// 	// 	filter(function(key) {
	// 	// 	  return criteria(object, key);
	// 	// 	}).
	// 	// 	reduce(function(acc, key) {
	// 	// 	  acc[key] = object[key];
	// 	// 	  return acc;
	// 	// 	},
	// 	//   {});
	return p
}

func queryMask(query Query) Pattern {
	fn := func(t Triple, key string) bool {
		log.Fatal(`return typeof triple[key] !== 'object';`)
		return false
	}
	return objectMask(fn, query)
}

func variablesMask(query Query) Pattern {
	fn := func(t Triple, key string) bool {
		log.Fatal(`return triple[key] instanceof Variable;`)
		return false
	}
	return objectMask(fn, query)
}

func genKeys(t Triple) [][]byte {
	keys := make([][]byte, 0, len(defs))

	for k := range defs {
		keys = append(keys, genKey(k, t))
	}

	return keys
}

func genKey(key string, triple Triple) []byte {
	if h == nil {
		var err error
		var k [16]byte
		copy(k[:], []byte("honeybadger"))
		h, err = siphash.New128(k[:])
		if err != nil {
			log.Fatal(err)
		}
	}

	result := make([]byte, len(key), 3+(3*16)) //3 letters of index + 3 128bit hashes
	copy(result[:], []byte(key))

	for _, d := range defs[key] {
		valid := false
		switch d {
		case "subject":
			if triple.Subject != nil {
				h.Write(triple.Subject)
				valid = true
			}
		case "predicate":
			if triple.Predicate != nil {
				h.Write(triple.Predicate)
				valid = true
			}
		case "object":
			if triple.Object != nil {
				h.Write(triple.Object)
				valid = true
			}
		default:
			runtime.Breakpoint()
		}

		if valid {
			subKey := h.Sum(nil)
			result = append(result, subKey...)
		}
		h.Reset()
	}

	return result
}

type streamType int

const (
	joinStream streamType = iota
)

//Query x
type Query struct {
	Prefix []byte
	Limit  uint
	stream streamType
}

func createQuery(pattern Pattern, options streamOptions) Query {
	types := typesFromPattern(pattern)
	preferiteIndex := options.Index
	index := findIndex(types, preferiteIndex)
	key := genKey(index, pattern.Triple)

	query := Query{
		Prefix: key,
		Limit:  pattern.Limit,
	}

	return query
}

func typesFromPattern(pattern Pattern) []string {
	results := make([]string, 0, 3)

	if len(pattern.Triple.Subject) > 0 {
		results = append(results, "subject")
	}

	if len(pattern.Triple.Predicate) > 0 {
		results = append(results, "predicate")
	}

	if len(pattern.Triple.Object) > 0 {
		results = append(results, "object")
	}

	return results
}

func findIndex(types []string, preferiteIndex string) string {
	result := possibleIndexes(types)
	there := false
	for _, i := range result {
		if i == preferiteIndex {
			there = true
			break
		}
	}

	if len(preferiteIndex) > 0 && there {
		return preferiteIndex
	}

	return result[0]
}

func possibleIndexes(types []string) []string {
	results := []string{}
	for key, def := range defs {

		matches := 0
		valid := false

		for _, e := range def {
			found := false
			for _, x := range types {
				if x == e {
					matches++
					found = true
					break
				}
			}

			matchedAll := matches == len(types)
			if found && !matchedAll {
				continue
			}

			if matchedAll {
				valid = true
				break
			}

			break
		}

		if valid {
			results = append(results, key)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})

	return results
}
