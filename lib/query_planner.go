package honeybadger

import (
	"log"
	"sort"

	"github.com/pkg/errors"
)

//NewQueryPlanner x
func NewQueryPlanner(hb *HoneyBadger, options VariableQueryOptions, queryPatterns ...QueryPattern) ([]QueryPattern, error) {
	results := queryPatterns //dupes?

	queriesSizes := make([]uint, len(queryPatterns))
	for i, qp := range queryPatterns {
		err := checkPatternForCollisions(qp)
		if err != nil {
			return nil, errors.Wrapf(err, "pattern collision in %+v", qp)
		}
		rdfCount, _ := hb.Count(qp)
		queriesSizes[i] = rdfCount

		log.Fatal("not implemented")
	}

	sort.Slice(queryPatterns, func(i, j int) bool {
		return queriesSizes[i] < queriesSizes[j]
	})

	if options.Algorithm == Sort && len(results) > 1 {
		var sum *QueryPattern

		for _, qp := range results {
			sum = doSortQueryPlan(sum, &qp)
		}
	}

	return results, nil
}

func checkPatternForCollisions(qp QueryPattern) error {
	s := qp.Triple.Subject
	p := qp.Triple.Predicate
	o := qp.Triple.Object
	sV := qp.Variables.Subject
	pV := qp.Variables.Predicate
	oV := qp.Variables.Object

	if s != nil && len(s) > 0 && sV != nil {
		return errors.New("uses both subject and bound variable")
	}

	if p != nil && len(p) > 0 && pV != nil {
		return errors.New("uses both predicate and bound variable")
	}

	if o != nil && len(o) > 0 && oV != nil {
		return errors.New("Pattern uses both object and bound variable")
	}

	return nil
}
func doSortQueryPlan(first, second *QueryPattern) *QueryPattern {
	// if first == nil || first.stream == joinStream {
	if first == nil {
		return nil
	}

	log.Fatal("not implemented")

	// firstQueryMask := queryMask(first)
	//   secondQueryMask := queryMask(second)
	//   firstVariablesMask := variablesMask(first)
	//   secondVariablesMask := variablesMask(second)

	//   firstVariables = Object.keys(firstVariablesMask).map(function(key) {
	// 	  return firstVariablesMask[key];
	// 	})

	//   , secondVariables = Object.keys(secondVariablesMask).map(function(key) {
	// 	  return secondVariablesMask[key];
	// 	})

	//   , variableKey = function(obj, variable) {
	// 	  return Object.keys(obj).filter(function(key) {
	// 		return obj[key].name === variable.name;
	// 	  })[0];
	// 	}

	//   , commonVariables = firstVariables.filter(function(firstVar) {
	// 	  return secondVariables.some(function(secondVar) {
	// 		return firstVar.name === secondVar.name;
	// 	  });
	// 	})

	//   , firstIndexArray = Object.keys(firstQueryMask)
	//   , secondIndexArray = Object.keys(secondQueryMask)

	//   , commonValueKeys = firstIndexArray.filter(function(key) {
	// 	  return secondIndexArray.indexOf(key) >= 0;
	// 	})

	//   , firstIndexes
	//   , secondIndexes;

	// if len(commonVariables) ==  0 {
	//   return nil
	// }

	// if first.stream == nil {
	// 	first.stream = JoinStream
	// }

	// firstIndexArray := firstIndexArray.filter(function(key) {
	//   return commonValueKeys.indexOf(key) < 0;
	// })

	// secondIndexArray := secondIndexArray.filter(function(key) {
	//   return commonValueKeys.indexOf(key) < 0;
	// })

	// for _,key := range commonValueKeys {
	//   firstIndexArray.unshift(key)
	//   secondIndexArray.unshift(key)
	// }

	// sort.Slice(commonVariables, func(i,j int) bool {
	// 	a:= commonVariables[i]
	// 	b:= commonVariables[j]
	// 	return a < b
	// })

	// for _,commonVar := range commonVariables {
	//   firstIndexArray = append( firstIndexArray,variableKey(firstVariablesMask, commonVar))
	//   secondIndexArray =append(secondIndexArray,variableKey(secondVariablesMask, commonVar))
	// })

	// firstIndexes = orderedPossibleIndex(firstIndexArray)
	// secondIndexes = orderedPossibleIndex(secondIndexArray)

	// first.index = first.index || firstIndexes[0]
	// second.index = secondIndexes[0]
	// second.stream = SortJoinStream

	return second
}
