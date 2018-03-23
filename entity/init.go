package entity

import (
	"dotacrawler/config"
)

type Entity struct {
	VpGame VpGame
}

var confVpGame config.VpGame

func NewEntity() Entity {
	confVpGame = config.Get().GetVpGame()
	return Entity{
		VpGame: newVpGame(),
	}
}
