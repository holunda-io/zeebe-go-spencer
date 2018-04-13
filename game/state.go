package game

type State struct {
	BaddieHealth int    `msgpack:"health"`
	Decision     string `msgpack:"attack"`
}

type PlayerSetting struct {
	NormalAttack    int
	SpecialAttack   int
	AdditionalRange int
}

func NewGame() *State {
	gameState := new(State)
	gameState.BaddieHealth = 100
	return gameState
}
