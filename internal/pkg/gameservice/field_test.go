package gameservice

import (
	"reflect"
	"testing"
)

func TestFieldCorrectTurn(t *testing.T) {
	type tcase []Point
	tCases := []tcase{
		[]Point{{X: 1, Y: 1}, {X: 3, Y: 1}},
		[]Point{{X: 1, Y: 1}, {X: 1, Y: 3}},
	}

	for _, c := range tCases {
		field := NewFieldModel(false)
		err := field.TryTurn(c, "1")
		if err != nil {
			t.Errorf("Unexpected err: %v", err)
		}
	}
}

func TestFieldIncorrectTurn(t *testing.T) {
	type tcase []Point
	tCases := []tcase{
		[]Point{{X: 1, Y: 0}, {X: 1, Y: 3}},
		[]Point{{X: 2, Y: 4}, {X: 4, Y: 4}},
		[]Point{{X: 1, Y: 1}, {X: 2, Y: 2}},
		[]Point{{X: -1, Y: 0}, {X: 0, Y: 0}},
		[]Point{{X: 0, Y: -1}, {X: 0, Y: 0}},
		[]Point{{X: 4, Y: 4}, {X: 4, Y: 5}},
		[]Point{{X: 4, Y: 4}, {X: 5, Y: 4}},
	}

	for _, c := range tCases {
		field := NewFieldModel(false)
		field.data = [][]string{
			{"", "*", "*", "*", "*"},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"2", "2", "2", "", ""},
		}
		field.freeCells = 18

		oldState := field.data
		err := field.TryTurn(c, "1")

		if err == nil {
			t.Errorf("Unexpected no err: %v", err)
		}
		if !reflect.DeepEqual(field.data, oldState) {
			t.Errorf("Incorrect turn was applied to real data")
		}
	}
}
