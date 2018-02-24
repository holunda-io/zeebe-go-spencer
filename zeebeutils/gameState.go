package zeebeutils

type GameState struct {
	BaddieHealth int    `msgpack:"health"`
	Decision     string `msgpack:"attack"`
}

type PlayerSetting struct {
	NormalAttack    int
	SpecialAttack   int
	AdditionalRange int
}