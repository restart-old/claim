package claim

import (
	"sync"

	"github.com/df-mc/dragonfly/server/player"
)

var players sync.Map

func StorePlayerClaim(p *player.Player, c Claim) {
	players.Store(p, c)
}
func LoadPlayerClaim(p *player.Player) (Claim, bool) {
	c, ok := players.Load(p)
	if ok {
		if c, ok := c.(Claim); ok {
			return c, true
		}
	}
	return nil, false
}
