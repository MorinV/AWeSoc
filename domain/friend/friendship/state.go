package friendship

import (
	"errors"
)

type State struct {
	Name  string
	Value int
}

var States = map[string]int{
	"pending":  1,
	"accepted": 2,
}

var (
	ErrNotExistingState = errors.New("не существующий статус дружбы")
)

func NewStateFromValue(stateValue int) (State, error) {
	for name, value := range States {
		if value == stateValue {
			return State{Value: value, Name: name}, nil
		}
	}
	return State{}, ErrNotExistingState
}

func NewStateFromName(stateName string) (State, error) {
	value, found := States[stateName]
	if found {
		return State{Name: stateName, Value: value}, nil
	}

	return State{}, ErrNotExistingState
}

func (s State) String() string {
	return s.Name
}

func (s State) Int() int {
	return s.Value
}
