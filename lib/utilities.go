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
	streamTypeJoin streamType = iota
)

//Query x
type Query struct {
	Prefix []byte
	Limit  uint
	stream streamType
}

func createQuery(qp QueryPattern, options streamOptions) Query {
	types := typesFromPattern(qp)
	preferiteIndex := options.Index
	index := findIndex(types, preferiteIndex)
	key := genKey(index, qp.Triple)

	query := Query{
		Prefix: key,
		Limit:  qp.Limit,
	}

	return query
}

func typesFromPattern(qp QueryPattern) []string {
	results := make([]string, 0, 3)

	if len(qp.Triple.Subject) > 0 {
		results = append(results, "subject")
	}

	if len(qp.Triple.Predicate) > 0 {
		results = append(results, "predicate")
	}

	if len(qp.Triple.Object) > 0 {
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

type maskUpdaterFn func(Solutions, []string) Solutions

func maskUpdater(qp QueryPattern) maskUpdaterFn {
	variables := variablesMask(qp)

	return func(solutions Solutions, mask []string) Solutions {
		newMask := Solutions{}

		for _, v := range variables {
			if v.isBound(solutions) {
				newMask[v.Name] = solutions[v.Name]
			}

			log.Fatal("Solutions have filters?")
			// newMask.Filter = qp.Filter
		}

		return newMask
		// 	return Object.keys(variables).reduce(function(newMask, key) {
		// 	var variable = variables[key];
		// 	if (variable.isBound(solution)) {
		// 	  newMask[key] = solution[variable.name];
		// 	}
		// 	newMask.filter = pattern.filter;
		// 	return newMask;
		//   }, Object.keys(mask).reduce(function(acc, key) {
		// 	acc[key] = mask[key];
		// 	return acc;
		//   }, {}));
	}
}

func variablesMask(qp QueryPattern) map[string]Variable {
	variables := map[string]Variable{}
	s := qp.Variables.Subject
	p := qp.Variables.Predicate
	o := qp.Variables.Object

	if s != nil {
		variables["subject"] = *s
	}
	if p != nil {
		variables["predicate"] = *p
	}
	if o != nil {
		variables["object"] = *o
	}
	return variables
}

type matcherFn func(Solutions, Triple)

func matcher(qp QueryPattern) matcherFn {
	queryVariables := variablesMask(qp)

	//Real matcher?
	return func(solutions Solutions, triple Triple) {
		for key, v := range queryVariables {
			var value []byte
			switch key {
			case "subject":
				value = triple.Subject
			case "predicate":
				value = triple.Predicate
			case "object":
				value = triple.Object
			default:
				log.Fatal("on noes")
			}
			v.bind(solutions, value)
		}
	}
}
