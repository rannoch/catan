package domain

type GameStatePlayerSelectingWhoToRob struct {
	game *Game

	GameStateDefault
}

func (GameStatePlayerSelectingWhoToRob) RobPlayer(playerColor Color, targetColor Color) error {
	panic("implement me")
}

func (GameStatePlayerSelectingWhoToRob) Apply(eventMessage EventMessage, isNew bool) {
	panic("implement me")
}
