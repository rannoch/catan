package grid

//           _ _
//         /     \
//    _ _ / (0,4) \ _ _
//  /     \       /     \
// / (0,3) \ _ _ / (1,5) \
// \       /     \       /
//  \ _ _ / (1,4) \ _ _ /
//  /     \       /     \
// / (1,3) \ _ _ / (2,5) \
// \       /     \       /
//  \ _ _ / (2,4) \ _ _ /
//        \       /
//         \ _ _ /
type HexagonGridWithOffsetCoordsCalculator struct{}

func (h HexagonGridWithOffsetCoordsCalculator) IntersectionAdjacentHexes(intersectionCoord IntersectionCoord) []HexCoord {
	if intersectionCoord.D == L {
		return []HexCoord{
			{R: intersectionCoord.R - 1, C: intersectionCoord.C - 1},
			{R: intersectionCoord.R, C: intersectionCoord.C},
			{R: intersectionCoord.R, C: intersectionCoord.C - 1},
		}
	} else if intersectionCoord.D == R {
		return []HexCoord{
			{R: intersectionCoord.R, C: intersectionCoord.C},
			{R: intersectionCoord.R, C: intersectionCoord.C + 1},
			{R: intersectionCoord.R + 1, C: intersectionCoord.C + 1},
		}
	}

	return nil
}

func (h HexagonGridWithOffsetCoordsCalculator) IntersectionAdjacentPaths(intersectionCoord IntersectionCoord) []PathCoord {
	if intersectionCoord.D == L {
		return []PathCoord{
			{R: intersectionCoord.R, C: intersectionCoord.C - 1, D: N},
			{R: intersectionCoord.R, C: intersectionCoord.C, D: W},
			{R: intersectionCoord.R, C: intersectionCoord.C - 1, D: E},
		}
	} else if intersectionCoord.D == R {
		return []PathCoord{
			{R: intersectionCoord.R, C: intersectionCoord.C, D: E},
			{R: intersectionCoord.R + 1, C: intersectionCoord.C, D: N},
			{R: intersectionCoord.R + 1, C: intersectionCoord.C, D: W},
		}
	}

	return nil
}

func (h HexagonGridWithOffsetCoordsCalculator) HexAdjacentIntersections(hexCoord HexCoord) []IntersectionCoord {
	return []IntersectionCoord{
		{R: hexCoord.R - 1, C: hexCoord.C - 1, D: R},
		{R: hexCoord.R, C: hexCoord.C + 1, D: L},
		{R: hexCoord.R, C: hexCoord.C, D: R},
		{R: hexCoord.R + 1, C: hexCoord.C + 1, D: L},
		{R: hexCoord.R, C: hexCoord.C - 1, D: R},
		{R: hexCoord.R, C: hexCoord.C, D: L},
	}
}

func (h HexagonGridWithOffsetCoordsCalculator) HexAdjacentPaths(hexCoord HexCoord) []PathCoord {
	return []PathCoord{
		{R: hexCoord.R, C: hexCoord.C, D: W},
		{R: hexCoord.R, C: hexCoord.C, D: N},
		{R: hexCoord.R, C: hexCoord.C, D: E},
		{R: hexCoord.R + 1, C: hexCoord.C + 1, D: W},
		{R: hexCoord.R + 1, C: hexCoord.C, D: N},
		{R: hexCoord.R, C: hexCoord.C - 1, D: E},
	}
}

type HexCoord struct {
	R int64 // row
	C int64 // column
}

type IntersectionCoord struct {
	R int64
	C int64
	D IntersectionDirection
}

type IntersectionDirection string

const (
	L IntersectionDirection = "left"
	R IntersectionDirection = "right"
)

type PathCoord struct {
	R int64
	C int64
	D PathDirection
}

type PathDirection string

const (
	W PathDirection = "west"
	N PathDirection = "north"
	E PathDirection = "east"
)