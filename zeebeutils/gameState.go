package zeebeutils

type GameState struct {
	Health   int 		`msgpack:"health"`
	Decision string 	`msgpack:"attack"`
}