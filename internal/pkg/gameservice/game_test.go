package gameservice

import (
	"testing"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

func TestGameModel_DoTurn(t *testing.T) {
	type fields struct {
		room      *Room
		field     *FieldModel
		players   []models.RoomPlayer
		whoseTurn int
		ready     bool
		running   bool
		//cellsCount int
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
				players: []models.RoomPlayer{
					{
						Username: "u1",
						Avatar:   "",
					},
					{
						Username: "u2",
						Avatar:   "",
					},
				},
				whoseTurn: 0,
				ready:     true,
				running:   true,
				field:     NewFieldModel(false),
			},
			args: args{
				a: GameMessage{
					GenericMessage{
						MType: "game",
						User:  "u1",
					},
					struct {
						Coords []Point `json:"coords"`
					}{
						Coords: []Point{{X: 1, Y: 1}, {X: 3, Y: 1}},
					},
				},
			},
			wantErr: false,
		},

		{
			name: "IncorrectTurn",
			fields: fields{
				players: []models.RoomPlayer{
					{
						Username: "u1",
						Avatar:   ""},
					{
						Username: "u2",
						Avatar:   "",
					},
				},
				whoseTurn: 1,
				ready:     true,
				running:   true,
				field:     NewFieldModel(false),
			},
			args: args{
				a: GameMessage{
					GenericMessage{
						MType: "game",
						User:  "u1",
					},
					struct {
						Coords []Point `json:"coords"`
					}{
						Coords: []Point{{X: 1, Y: 1}, {X: 3, Y: 3}},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameModel{
				room:      tt.fields.room,
				field:     tt.fields.field,
				players:   tt.fields.players,
				whoseTurn: tt.fields.whoseTurn,
				ready:     tt.fields.ready,
				running:   tt.fields.running,
			}
			if err := g.DoTurn(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("GameModel.DoTurn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
