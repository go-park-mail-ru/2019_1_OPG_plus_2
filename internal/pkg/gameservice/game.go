package gameservice

import (
    "errors"

    "2019_1_OPG_plus_2/internal/pkg/randomgenerator"
)

const (
    freeFreq   int    = 3
    height     int    = 5
    width      int    = 5
    freeCell   string = " "
    lockedCell string = "*"
)

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

type FieldModel struct {
    data      [][]string
    width     int
    height    int
    freeCells int
}

func NewFieldModel(withBlocked bool) *FieldModel {
    freeCells := 0
    field := make([][]string, height)
    for r := range field {
        field[r] = make([]string, width)
        for c := range field[r] {
            if withBlocked {
                switch randomgenerator.RandomInt(0, freeFreq) {
                case 0:
                    field[r][c] = lockedCell
                default:
                    field[r][c] = freeCell
                    freeCells++
                }
            } else {
                field[r][c] = freeCell
                freeCells++
            }
        }
    }

    return &FieldModel{
        data:      field,
        width:     width,
        height:    height,
        freeCells: freeCells,
    }
}

func (f *FieldModel) TryTurn(coords []Point, char string) (err error) {
    if len(coords) != 2 {
        return errors.New("incorrect count of points")
    }
    if coords[0].X < 0 || f.width <= coords[0].X || coords[1].X < 0 || f.width <= coords[1].X ||
        coords[0].Y < 0 || f.height <= coords[0].Y || coords[1].Y < 0 || f.height <= coords[1].Y {
        return errors.New("out of field")
    }
    
    minX := min(coords[0].X, coords[1].X)
    maxX := max(coords[0].X, coords[1].X)
    minY := min(coords[0].Y, coords[1].Y)
    maxY := max(coords[0].Y, coords[1].Y)
    if maxX-minX != 0 && maxY-minY != 0 {
        return errors.New("not row or column")
    }

    // Check for blocked cells
    for x := minX; x <= maxX; x++ {
        for y := minY; y <= maxY; y++ {
            if f.data[y][x] != freeCell {
                return errors.New("some cells is blocked")
            }
        }
    }

    // Make turn
    for x := minX; x <= maxX; x++ {
        for y := minY; y <= maxY; y++ {
            f.data[y][x] = char
        }
    }
    f.freeCells -= maxX - minX + maxY - minY
    return nil
}

type GameModel struct {
    room *Room

    field     *FieldModel
    players   []string
    whoseTurn int

    ready   bool
    running bool
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
    g.whoseTurn = (g.whoseTurn + 1) % len(g.players)
    return nil
}

func (g *GameModel) Check() bool {
    return g.field.freeCells <= 0
}

func (g *GameModel) IsReady() bool {
    return g.ready
}

func (g *GameModel) IsRunning() bool {
    return g.running
}

func (g *GameModel) GetField() [][]string {
    return g.field.data
}

func (g *GameModel) Init() {
    g.whoseTurn = randomgenerator.RandomInt(0, len(g.players))
    g.field = NewFieldModel(true)
}
