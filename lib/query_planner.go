package honeybadger

import (
	"log"
	"sort"
)

//QueryPlanner  x
type QueryPlanner []Query

//NewQueryPlanner x
func NewQueryPlanner(hb *HoneyBadger, options VariableQueryOptions, queries ...Query) QueryPlanner {
	results := queries //dupes?

	queriesSizes := make([]uint, len(queries))
	for i, q := range queries {
		newQ := queryMask(q)
		// rang := createQuery(newQ, streamOptions{})

		rdfCount, _ := hb.Count(newQ)
		queriesSizes[i] = rdfCount

		log.Fatal("not implemented")
	}

	sort.Slice(queries, func(i, j int) bool {
		return queriesSizes[i] < queriesSizes[j]
	})

	if options.Algorithm == Sort && len(results) > 1 {
		var sum *Query

		for _, q := range results {
			sum = doSortQueryPlan(sum, &q)
		}
	}

	return results
}

func doSortQueryPlan(first, second *Query) *Query {
	if first == nil || first.stream == joinStream {
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
