package honeybadger

// Solution x
type Solution map[string]string

//Variable x
type Variable struct {
	Name string
}

func (v *Variable) bind(solution Solution, value string) Solution {

	if !v.isBindable(solution, value) {
		return nil
	}

	newSolution := Solution{}
	for k, v := range solution {
		newSolution[k] = v
	}
	newSolution[v.Name] = value
	return newSolution
}

func (v *Variable) isBound(solution Solution) bool {
	_, ok := solution[v.Name]
	return ok
}

func (v *Variable) isBindable(solution Solution, value string) bool {
	currentValue, ok := solution[v.Name]
	return !ok || currentValue == value
}
