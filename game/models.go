package game

// GameInfo is send on game start
type GameInfo struct {
	maxTurns            int
	players             []string
	mapWidth, mapHeight int
}

type TurnInfo struct {
	turn int
}

type PlayerTurn struct {
	PlayerId int
}

type PlayerInfo struct {
	name string
}
