package game

// GameInfo is send on game start
type gameInfo struct {
	MaxTurns  int      `json:"maxTurns"`
	Players   []string `json:"players"`
	MapWidth  int      `json:"mapWidth"`
	MapHeight int      `json:"mapHeight"`
}

type gameResult struct {
	Winner string `json:"winner"`
}

type turnInfo struct {
	Turn int `json:"turn"`
}

type playerTurn struct {
	PlayerID int
}
