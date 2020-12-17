package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGame_NextTurnColor(t *testing.T) {
	type fields struct {
		Status     GameStatus
		TurnOrder  []Color
		TotalTurns int64
	}
	tests := []struct {
		name   string
		fields fields
		want   Color
	}{
		{
			name: "set-up phase forward",
			fields: fields{
				Status:     GameStatusInitialSetup,
				TurnOrder:  []Color{Red, Blue, Green, Yellow},
				TotalTurns: 0,
			},
			want: Blue,
		},
		{
			name: "set-up phase backward",
			fields: fields{
				Status:     GameStatusInitialSetup,
				TurnOrder:  []Color{Red, Blue, Green, Yellow},
				TotalTurns: 3,
			},
			want: Yellow,
		},
		{
			name: "play phase turn after 12",
			fields: fields{
				Status:     GameStatusPlay,
				TurnOrder:  []Color{Red, Blue, Green, Yellow},
				TotalTurns: 12,
			},
			want: Blue,
		},
		{
			name: "play phase turn after 3",
			fields: fields{
				Status:     GameStatusPlay,
				TurnOrder:  []Color{Red, Blue, Green, Yellow},
				TotalTurns: 3,
			},
			want: Red,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Game{
				status:     tt.fields.Status,
				turnOrder:  tt.fields.TurnOrder,
				totalTurns: tt.fields.TotalTurns,
			}
			got := g.NextTurnColor()

			assert.Equal(t, tt.want, got)
		})
	}
}
