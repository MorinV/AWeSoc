package person

import "errors"

type Gender struct {
	Name  string
	Value int
}

var Genders = map[string]int{
	"Мужчина": 1,
	"Женщина": 2,
}

var (
	ErrNotExistingGender = errors.New("не существующий гендер")
)

func NewGenderFromValue(stateValue int) (Gender, error) {
	for name, value := range Genders {
		if value == stateValue {
			return Gender{Value: value, Name: name}, nil
		}
	}
	return Gender{}, ErrNotExistingGender
}

func NewGenderFromName(stateName string) (Gender, error) {
	value, found := Genders[stateName]
	if found {
		return Gender{Name: stateName, Value: value}, nil
	}

	return Gender{}, ErrNotExistingGender
}

func (g Gender) String() string {
	return g.Name
}

func (g Gender) Int() int {
	return g.Value
}
