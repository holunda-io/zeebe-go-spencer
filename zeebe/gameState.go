package zeebe

type GameState struct {
	BaddieHealth int    `msgpack:"health"`
	Decision     string `msgpack:"attack"`
}

type PlayerSetting struct {
	NormalAttack    int
	SpecialAttack   int
	AdditionalRange int
}

func newGame() *GameState {
	gameState := new(GameState)
	gameState.BaddieHealth = 100
	return gameState
}