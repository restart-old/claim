package claim

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"log"
)

type Handler interface {
	HandleEnterClaim(ctx *event.Context, p *player.Player)
	HandleLeaveClaim(ctx *event.Context, p *player.Player)

	HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack)
	HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64)
}

type NopHandler struct{}

func (NopHandler) HandleEnterClaim(ctx *event.Context, p *player.Player)                         {}
func (NopHandler) HandleLeaveClaim(ctx *event.Context, p *player.Player)                         {}
func (NopHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack)        {}
func (NopHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64) {}

// ClaimHandler is the handler which is used to handle:
// When a block is broken in a claim.
// When a player enters or leaves a claim.
// And when a player is hurt in a claim.
type ClaimHandler struct {
	player.NopHandler
	p *player.Player
	Loader
}

func (*ClaimHandler) Name() string { return "ClaimHandler" }

// NewClaimHandler returns a new *ClaimHandler.
func NewClaimHandler(p *player.Player, loader Loader) *ClaimHandler {
	return &ClaimHandler{
		p:      p,
		Loader: loader,
	}
}

// HandleBlockBreak handles when a block is broken,
// and cancels the event if breaking blocks are not allowed in the claim they were broken in.
func (c *ClaimHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	claim, err := c.LoadWithPos(pos.Vec3())
	if err != nil {
		log.Println(err)
		return
	}
	claim.h.HandleBlockBreak(ctx, pos, drops)
}

func (c *ClaimHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64) {
	claim, err := c.LoadWithPos(e.Position())
	if err != nil {
		log.Println(err)
		return
	}
	claim.h.HandleAttackEntity(ctx, e, force, height)
}

func actuallyMovedXZ(old, new mgl64.Vec3) bool {
	return old.X() != new.X() || old.Z() != new.Z()
}

func (c *ClaimHandler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if actuallyMovedXZ(c.p.Position(), newPos) {
		claim, err := c.LoadWithPos(newPos)
		if err != nil {
			log.Println(err)
			return
		}
		claim.Enter(ctx, c.p)
	}
}
