package claim

import (
	"sync"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dragonfly-on-steroids/area"
)

var players sync.Map

var claims []*Claim

func Store(c *Claim) {
	claims = append(claims, c)
}
func Delete(c *Claim) {
	for n, claim := range claims {
		if claim == c {
			claims = append(claims[:n], claims[1+n:]...)
		}
	}
}

func Claims() []*Claim { return claims }

func NewClaim(name string, w *world.World, area area.Vec2) *Claim {
	return &Claim{
		world: w,
		area:  area,
		h:     NopHandler{},
		name:  name,
	}
}

type Claim struct {
	name   string
	world  *world.World
	area   area.Vec2
	hMutex sync.RWMutex
	h      Handler
}

func (c *Claim) Name() string        { return c.name }
func (c *Claim) World() *world.World { return c.world }
func (c *Claim) Area() area.Vec2     { return c.area }
func (c *Claim) handler() Handler    { return c.h }
func (c *Claim) Handle(h Handler) {
	c.hMutex.Lock()
	defer c.hMutex.Unlock()
	if h == nil {
		h = NopHandler{}
	}
	c.h = h
}
func (c *Claim) EnterClaim(ctx *event.Context, p *player.Player) {
	if claim, _ := players.Load(p); claim != c {
		if claim == nil {
			Wilderness.h.HandleLeaveClaim(ctx, p)
		}
		c.h.HandleEnterClaim(ctx, p)
		ctx.Continue(func() {
			players.Store(p, c)
		})
	}
}
func (c *Claim) LeaveClaim(ctx *event.Context, p *player.Player) {
	if claim, _ := players.Load(p); claim == c {
		c.h.HandleLeaveClaim(ctx, p)
		ctx.Continue(func() {
			players.Delete(p)
		})
		if claim != Wilderness {
			Wilderness.h.HandleEnterClaim(ctx, p)
		}
	}
}
