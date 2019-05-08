package gameservice

import "testing"

func TestGameModel_DoTurn(t *testing.T) {
	type fields struct {
		room       *Room
		field      *FieldModel
		players    []string
		whoseTurn  int
		ready      bool
		running    bool
		cellsCount int
	}
	type args struct {
		a GameMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "CorrectTurn",
			fields: fields{
				players:   []string{"u1", "u2"},
				whoseTurn: 0,
				ready:     true,
				running:   true,
				field:     NewFieldModel(),
			},
			args: args{
				a: GameMessage{
					GenericMessage{
						MType: "game",
						User:  "u1",
					},
					struct {
						Coords []int `json:"coords"`
					}{
						Coords: []int{1, 2, 3, 4},
					},
				},
			},
			wantErr: false,
		},

		{
			name: "IncorrectTurn",
			fields: fields{
				players:   []string{"u1", "u2"},
				whoseTurn: 1,
				ready:     true,
				running:   true,
				field:     NewFieldModel(),
			},
			args: args{
				a: GameMessage{
					GenericMessage{
						MType: "game",
						User:  "u1",
					},
					struct {
						Coords []int `json:"coords"`
					}{
						Coords: []int{23, 24, 25},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameModel{
				room:       tt.fields.room,
				field:      tt.fields.field,
				players:    tt.fields.players,
				whoseTurn:  tt.fields.whoseTurn,
				ready:      tt.fields.ready,
				running:    tt.fields.running,
				cellsCount: tt.fields.cellsCount,
			}
			if err := g.DoTurn(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("GameModel.DoTurn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
