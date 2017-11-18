package honeybadger

// Solutions x
type Solutions map[string]string

//Variable x
type Variable struct {
	Name string
}

func (v *Variable) bind(solutions Solutions, value string) Solutions {

	if !v.isBindable(solutions, value) {
		return nil
	}

	newSolutions := Solutions{}
	for k, v := range solutions {
		newSolutions[k] = v
	}
	newSolutions[v.Name] = value
	return newSolutions
}

func (v *Variable) isBound(solutions Solutions) bool {
	_, ok := solutions[v.Name]
	return ok
}

func (v *Variable) isBindable(solutions Solutions, value string) bool {
	currentValue, ok := solutions[v.Name]
	return !ok || currentValue == value
}
