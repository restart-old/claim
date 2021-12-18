package claim

import (
	"sync"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dragonfly-on-steroids/area"
)

var players sync.Map

func NewClaim(name string, area area.Vec2) *Claim {
	return &Claim{
		area: area,
		h:    NopHandler{},
		name: name,
	}
}

type Claim struct {
	name   string
	world  *world.World
	area   area.Vec2
	hMutex sync.RWMutex
	h      Handler
}

func (c *Claim) Compare(claim2 interface{}) bool {
	if claim, ok := claim2.(*Claim); ok {
		return c.name == claim.name
	}
	return false
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
func (c *Claim) Enter(ctx *event.Context, p *player.Player) {
	if claim, _ := players.Load(p); !c.Compare(claim) {
		c.h.HandleEnterClaim(ctx, p)
		ctx.Continue(func() {
			if claim, ok := claim.(*Claim); ok {
				claim.Leave(ctx, p)
			}
			players.Store(p, c)
		})
	}
}
func (c *Claim) Leave(ctx *event.Context, p *player.Player) {
	if claim, _ := players.Load(p); c.Compare(claim) {
		c.h.HandleLeaveClaim(ctx, p)
		ctx.Continue(func() {
			players.Delete(p)
		})
	}
}
