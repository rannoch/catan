package grid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHexagonGridWithOffsetCoords_VertexAdjacentHexes(t *testing.T) {
	type args struct {
		intersectionCoord IntersectionCoord
	}
	tests := []struct {
		name string
		args args
		want []HexCoord
	}{
		{
			name: "looking hexes for intersection (1,2,R), should be (1,2)(1,3)(2,3)",
			args: args{
				IntersectionCoord{
					R: 1,
					C: 2,
					D: R,
				},
			},
			want: []HexCoord{
				{R: 1, C: 2}, {R: 1, C: 3}, {R: 2, C: 3},
			},
		},
		{
			name: "looking hexes for intersection (1,3,L), should be (0,2)(1,3)(1,2)",
			args: args{
				IntersectionCoord{
					R: 1,
					C: 3,
					D: L,
				},
			},
			want: []HexCoord{
				{R: 0, C: 2}, {R: 1, C: 3}, {R: 1, C: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HexagonGridWithOffsetCoordsCalculator{}

			got := h.IntersectionAdjacentHexes(tt.args.intersectionCoord)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHexagonGridWithOffsetCoords_HexAdjacentVertices(t *testing.T) {
	type args struct {
		hexCoord HexCoord
	}
	tests := []struct {
		name string
		args args
		want []IntersectionCoord
	}{
		{
			name: "looking intersections for hex (1,3), should be (0,2,R),(1,4,L),(1,3,R),(2,4,L),(1,2,R),(1,3,L)",
			args: args{
				hexCoord: HexCoord{R: 1, C: 3},
			},
			want: []IntersectionCoord{
				{R: 0, C: 2, D: R},
				{R: 1, C: 4, D: L},
				{R: 1, C: 3, D: R},
				{R: 2, C: 4, D: L},
				{R: 1, C: 2, D: R},
				{R: 1, C: 3, D: L},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			h := HexagonGridWithOffsetCoordsCalculator{}

			got := h.HexAdjacentIntersections(tt.args.hexCoord)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHexagonGridWithOffsetCoords_VertexAdjacentEdges(t *testing.T) {
	type args struct {
		intersectionCoord IntersectionCoord
	}
	tests := []struct {
		name string
		args args
		want []PathCoord
	}{
		{
			name: "looking adjacent paths for intersection (1,3,L), should be (1,2,N),(1,3,W),(1,2,E)",
			args: args{
				intersectionCoord: IntersectionCoord{R: 1, C: 3, D: L},
			},
			want: []PathCoord{
				{R: 1, C: 2, D: N},
				{R: 1, C: 3, D: W},
				{R: 1, C: 2, D: E},
			},
		},
		{
			name: "",
			args: args{
				intersectionCoord: IntersectionCoord{R: 0, C: 3, D: R},
			},
			want: []PathCoord{
				{R: 0, C: 3, D: E},
				{R: 1, C: 3, D: N},
				{R: 1, C: 3, D: W},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			h := HexagonGridWithOffsetCoordsCalculator{}

			got := h.IntersectionAdjacentPaths(tt.args.intersectionCoord)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHexagonGridWithOffsetCoordsCalculator_HexAdjacentPaths(t *testing.T) {
	type args struct {
		hexCoord HexCoord
	}
	tests := []struct {
		name string
		args args
		want []PathCoord
	}{
		{
			name: "looking adjacent paths for hex (1,2)",
			args: args{
				hexCoord: HexCoord{R: 1, C: 2},
			},
			want: []PathCoord{
				{R: 1, C: 2, D: W},
				{R: 1, C: 2, D: N},
				{R: 1, C: 2, D: E},
				{R: 2, C: 3, D: W},
				{R: 2, C: 2, D: N},
				{R: 1, C: 1, D: E},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			h := HexagonGridWithOffsetCoordsCalculator{}
			got := h.HexAdjacentPaths(tt.args.hexCoord)

			assert.Equal(t, tt.want, got)
		})
	}
}
