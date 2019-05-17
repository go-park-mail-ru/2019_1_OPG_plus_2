package gameservice

import (
	"fmt"
	"math/rand"
)

type FieldModel struct {
	field     [][]string
	fieldCopy [][]string
}

func NewFieldModel() *FieldModel {
	return &FieldModel{
		field: [][]string{
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
		},
		fieldCopy: [][]string{
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
		},
	}
}

func (f *FieldModel) TryTurn(coords []int, char string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invaild turn")
		}
	}()
	copy(f.fieldCopy, f.field)
	for _, coord := range coords {
		x := coord / 5
		y := coord % 5
		if f.fieldCopy[x][y] != "" {
			return fmt.Errorf("invalid turn")
		}
		f.fieldCopy[x][y] = char
	}
	copy(f.field, f.fieldCopy)
	return nil
}

type GameModel struct {
	room *Room

	field     *FieldModel
	players   []string
	whoseTurn int

	ready      bool
	running    bool
	cellsCount int
}

func NewGameModel(room *Room) *GameModel {
	return &GameModel{
		room: room,
	}
}

func (g *GameModel) DoTurn(a GameMessage) error {
	err := g.field.TryTurn(a.Data.Coords, a.User)
	if err != nil {
		return err
	}
	g.cellsCount -= len(a.Data.Coords)
	g.whoseTurn = (g.whoseTurn + 1) % len(g.players)
	return nil
}

func (g *GameModel) Check() bool {
	return g.cellsCount <= 24
}

func (g *GameModel) IsReady() bool {
	return g.ready
}

func (g *GameModel) IsRunning() bool {
	return g.running
}

func (g *GameModel) Init() {
	g.whoseTurn = rand.Intn(len(g.players))
	g.field = NewFieldModel()
	g.cellsCount = 25
}
