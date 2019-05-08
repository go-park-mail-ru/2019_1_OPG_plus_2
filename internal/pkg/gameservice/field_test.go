package gameservice

import (
	"reflect"
	"testing"
)

func TestFieldCorrectTurn(t *testing.T) {
	type tcase []int
	tCases := []tcase{
		[]int{1, 2, 3},
		[]int{1, 6, 11},
	}

	for _, c := range tCases {
		field := NewFieldModel()
		err := field.TryTurn(c, "1")
		if err != nil {
			t.Errorf("Unexpected err: %v", err)
		}
	}
}

func TestFieldIncorrectTurn(t *testing.T) {
	type tcase []int
	tCases := []tcase{
		[]int{0, 1, 2},
		[]int{15, 20},
	}

	for _, c := range tCases {
		field := NewFieldModel()
		field.field = [][]string{
			{"", "*", "*", "*", "*"},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"2", "2", "2", "", ""},
		}
		field.fieldCopy = [][]string{
			{"", "*", "*", "*", "*"},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"2", "2", "2", "", ""},
		}

		oldState := field.field
		oldStateCopy := field.fieldCopy
		err := field.TryTurn(c, "1")

		if err == nil {
			t.Errorf("Unexpected no err: %v", err)
		}
		if !reflect.DeepEqual(field.field, oldState) {
			t.Errorf("Incorrect turn was applied to real field")
		}
		if !reflect.DeepEqual(field.fieldCopy, oldStateCopy) {
			t.Errorf("Incorrect turn was applied to field temp copy")
		}
	}
}
