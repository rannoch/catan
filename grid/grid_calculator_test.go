package grid

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
			name: "looking adjacent paths for intersection (0,0,R), should be (0,0,E),(1,1,N),(1,1,W)",
			args: args{
				intersectionCoord: IntersectionCoord{R: 0, C: 0, D: R},
			},
			want: []PathCoord{
				{R: 0, C: 0, D: E},
				{R: 1, C: 1, D: N},
				{R: 1, C: 1, D: W},
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

func TestHexagonGridWithOffsetCoordsCalculator_PathAdjacentIntersections(t *testing.T) {
	type args struct {
		pathCoord PathCoord
	}
	tests := []struct {
		name string
		args args
		want []IntersectionCoord
	}{
		{
			name: "looking adjacent intersections for path (1,2,E)",
			args: args{
				PathCoord{R: 1, C: 2, D: E},
			},
			want: []IntersectionCoord{
				{R: 1, C: 3, D: L},
				{R: 1, C: 2, D: R},
			},
		},
		{
			name: "looking adjacent intersections for path (1,2,N)",
			args: args{
				PathCoord{R: 1, C: 2, D: N},
			},
			want: []IntersectionCoord{
				{R: 0, C: 1, D: R},
				{R: 1, C: 3, D: L},
			},
		},
		{
			name: "looking adjacent intersections for path (1,2,W)",
			args: args{
				PathCoord{R: 1, C: 2, D: W},
			},
			want: []IntersectionCoord{
				{R: 1, C: 2, D: L},
				{R: 0, C: 1, D: R},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			h := HexagonGridWithOffsetCoordsCalculator{}
			if got := h.PathAdjacentIntersections(tt.args.pathCoord); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PathAdjacentIntersections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexagonGridWithOffsetCoordsCalculator_PathAdjacentPaths(t *testing.T) {
	type args struct {
		pathCoord PathCoord
	}
	tests := []struct {
		name string
		args args
		want []PathCoord
	}{
		{
			name: "looking adjacent paths for path (1,2,E)",
			args: args{
				PathCoord{R: 1, C: 2, D: E},
			},
			want: []PathCoord{
				{R: 1, C: 2, D: N},
				{R: 1, C: 3, D: W},
				{R: 2, C: 3, D: N},
				{R: 2, C: 3, D: W},
			},
		},
		{
			name: "looking adjacent paths for path (1,2,N)",
			args: args{
				PathCoord{R: 1, C: 2, D: N},
			},
			want: []PathCoord{
				{R: 1, C: 2, D: W},
				{R: 0, C: 1, D: E},
				{R: 1, C: 3, D: W},
				{R: 1, C: 2, D: E},
			},
		},
		{
			name: "looking adjacent paths for path (1,2,W)",
			args: args{
				PathCoord{R: 1, C: 2, D: W},
			},
			want: []PathCoord{
				{R: 0, C: 1, D: E},
				{R: 1, C: 2, D: N},
				{R: 1, C: 1, D: E},
				{R: 1, C: 1, D: N},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			h := HexagonGridWithOffsetCoordsCalculator{}
			if got := h.PathAdjacentPaths(tt.args.pathCoord); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PathAdjacentPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexagonGridWithOffsetCoordsCalculator_PathsJointIntersection(t *testing.T) {
	type args struct {
		pathCoord1 PathCoord
		pathCoord2 PathCoord
	}
	type want struct {
		intersectionCoord IntersectionCoord
		found             bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "looking joint intersection for paths (1,2,N)(1,2,E)",
			args: args{
				pathCoord1: PathCoord{R: 1, C: 2, D: N},
				pathCoord2: PathCoord{R: 1, C: 2, D: E},
			},
			want: want{
				intersectionCoord: IntersectionCoord{R: 1, C: 3, D: L},
				found:             true,
			},
		},
		{
			name: "looking joint intersection for paths (1,2,E)(2,3,W)",
			args: args{
				pathCoord1: PathCoord{R: 1, C: 2, D: E},
				pathCoord2: PathCoord{R: 2, C: 3, D: W},
			},
			want: want{
				intersectionCoord: IntersectionCoord{R: 1, C: 2, D: R},
				found:             true,
			},
		},
		{
			name: "looking joint intersection for paths (1,2,E)(2,3,W)",
			args: args{
				pathCoord1: PathCoord{R: 1, C: 2, D: E},
				pathCoord2: PathCoord{R: 100, C: 100, D: W},
			},
			want: want{
				intersectionCoord: IntersectionCoord{},
				found:             false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			h := HexagonGridWithOffsetCoordsCalculator{}
			intersection, found := h.PathsJointIntersection(tt.args.pathCoord1, tt.args.pathCoord2)

			assert.Equal(t, tt.want.intersectionCoord, intersection)
			assert.Equal(t, tt.want.found, found)
		})
	}
}

func TestHexagonGridWithOffsetCoordsCalculator_IntersectionAdjacentIntersections(t *testing.T) {
	type args struct {
		intersectionCoord IntersectionCoord
	}
	tests := []struct {
		name string
		args args
		want []IntersectionCoord
	}{
		{
			name: "R",
			args: args{
				intersectionCoord: IntersectionCoord{R: 0, C: 0, D: R},
			},
			want: []IntersectionCoord{
				{R: 0, C: 1, D: L},
				{R: 1, C: 2, D: L},
				{R: 1, C: 1, D: L},
			},
		},
		{
			name: "L",
			args: args{
				intersectionCoord: IntersectionCoord{R: 1, C: 1, D: L},
			},
			want: []IntersectionCoord{
				{R: 0, C: -1, D: R},
				{R: 0, C: 0, D: R},
				{R: 1, C: 0, D: R},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			h := HexagonGridWithOffsetCoordsCalculator{}
			if got := h.IntersectionAdjacentIntersections(tt.args.intersectionCoord); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntersectionAdjacentIntersections() = %v, want %v", got, tt.want)
			}
		})
	}
}
