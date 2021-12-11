package claim

import (
	"github.com/df-mc/dragonfly/server/entity/physics"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type wilderness struct {
	NopClaim
	w                  *world.World
	name               string
	enterMSG, leaveMSG string
}

func NewWilderness(w *world.World, enterMSG, leaveMSG string) wilderness {
	return wilderness{
		w:        w,
		name:     "Wilderness",
		enterMSG: enterMSG,
		leaveMSG: leaveMSG,
	}
}

func (w wilderness) AABB() physics.AABB  { return physics.AABB{} }
func (w wilderness) Name() string        { return w.name }
func (w wilderness) World() *world.World { return w.w }

func (w wilderness) Enter(p *player.Player) {
	p.Messagef(w.enterMSG)
}
func (w wilderness) Leave(p *player.Player) {
	p.Messagef(w.leaveMSG)
}
