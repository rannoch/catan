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
type HexagonGridCalculator interface {
	IntersectionAdjacentHexes(intersectionCoord IntersectionCoord) []HexCoord

	IntersectionAdjacentPaths(intersectionCoord IntersectionCoord) []PathCoord

	IntersectionAdjacentIntersections(intersectionCoord IntersectionCoord) []IntersectionCoord

	HexAdjacentIntersections(hexCoord HexCoord) []IntersectionCoord

	HexAdjacentPaths(hexCoord HexCoord) []PathCoord

	PathAdjacentIntersections(pathCoord PathCoord) []IntersectionCoord

	PathAdjacentPaths(pathCoord PathCoord) []PathCoord

	PathsJointIntersection(pathCoord1, pathCoord2 PathCoord) (IntersectionCoord, bool)
}

type HexagonGridWithOffsetCoordsCalculator struct{}

var _ HexagonGridCalculator = HexagonGridWithOffsetCoordsCalculator{}

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
			{R: intersectionCoord.R + 1, C: intersectionCoord.C + 1, D: N},
			{R: intersectionCoord.R + 1, C: intersectionCoord.C + 1, D: W},
		}
	}

	return nil
}

func (h HexagonGridWithOffsetCoordsCalculator) IntersectionAdjacentIntersections(intersectionCoord IntersectionCoord) []IntersectionCoord {
	if intersectionCoord.D == R {
		return []IntersectionCoord{
			{R: intersectionCoord.R, C: intersectionCoord.C + 1, D: L},
			{R: intersectionCoord.R + 1, C: intersectionCoord.C + 2, D: L},
			{R: intersectionCoord.R + 1, C: intersectionCoord.C + 1, D: L},
		}
	} else if intersectionCoord.D == L {
		return []IntersectionCoord{
			{R: intersectionCoord.R - 1, C: intersectionCoord.C - 2, D: R},
			{R: intersectionCoord.R - 1, C: intersectionCoord.C - 1, D: R},
			{R: intersectionCoord.R, C: intersectionCoord.C - 1, D: R},
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

func (h HexagonGridWithOffsetCoordsCalculator) PathAdjacentIntersections(pathCoord PathCoord) []IntersectionCoord {
	switch pathCoord.D {
	case W:
		return []IntersectionCoord{
			{R: pathCoord.R, C: pathCoord.C, D: L},
			{R: pathCoord.R - 1, C: pathCoord.C - 1, D: R},
		}
	case E:
		return []IntersectionCoord{
			{R: pathCoord.R, C: pathCoord.C + 1, D: L},
			{R: pathCoord.R, C: pathCoord.C, D: R},
		}
	case N:
		return []IntersectionCoord{
			{R: pathCoord.R - 1, C: pathCoord.C - 1, D: R},
			{R: pathCoord.R, C: pathCoord.C + 1, D: L},
		}
	}

	return nil
}

func (h HexagonGridWithOffsetCoordsCalculator) PathAdjacentPaths(pathCoord PathCoord) []PathCoord {
	switch pathCoord.D {
	case W:
		return []PathCoord{
			{R: pathCoord.R - 1, C: pathCoord.C - 1, D: E},
			{R: pathCoord.R, C: pathCoord.C, D: N},
			{R: pathCoord.R, C: pathCoord.C - 1, D: E},
			{R: pathCoord.R, C: pathCoord.C - 1, D: N},
		}
	case E:
		return []PathCoord{
			{R: pathCoord.R, C: pathCoord.C, D: N},
			{R: pathCoord.R, C: pathCoord.C + 1, D: W},
			{R: pathCoord.R + 1, C: pathCoord.C + 1, D: N},
			{R: pathCoord.R + 1, C: pathCoord.C + 1, D: W},
		}
	case N:
		return []PathCoord{
			{R: pathCoord.R, C: pathCoord.C, D: W},
			{R: pathCoord.R - 1, C: pathCoord.C - 1, D: E},
			{R: pathCoord.R, C: pathCoord.C + 1, D: W},
			{R: pathCoord.R, C: pathCoord.C, D: E},
		}
	}

	return nil
}

func (h HexagonGridWithOffsetCoordsCalculator) PathsJointIntersection(pathCoord1, pathCoord2 PathCoord) (IntersectionCoord, bool) {
	intersections1 := h.PathAdjacentIntersections(pathCoord1)
	intersections2 := h.PathAdjacentIntersections(pathCoord2)

	for _, intersectionCoord1 := range intersections1 {
		for _, intersectionCoord2 := range intersections2 {
			if intersectionCoord1 == intersectionCoord2 {
				return intersectionCoord1, true
			}
		}
	}

	return IntersectionCoord{}, false
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
	// W - west path direction
	W PathDirection = "west"
	// N - north path direction
	N PathDirection = "north"
	// E - east path direction
	E PathDirection = "east"
)
