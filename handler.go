package claim

import (
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
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
}

func (*ClaimHandler) Name() string { return "ClaimHandler" }

// NewClaimHandler returns a new *ClaimHandler.
func NewClaimHandler(p *player.Player) *ClaimHandler {
	return &ClaimHandler{
		p: p,
	}
}

// HandleBlockBreak handles when a block is broken,
// and cancels the event if breaking blocks are not allowed in the claim they were broken in.
func (c *ClaimHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	for _, claim := range Claims() {
		if claim.area.Vec2Within(mgl64.Vec2{pos.Vec3()[0], pos.Vec3()[2]}) {
			claim.h.HandleBlockBreak(ctx, pos, drops)
			return
		}
	}
	Wilderness.h.HandleBlockBreak(ctx, pos, drops)
}

func (c *ClaimHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64) {
	for _, claim := range Claims() {
		if claim.area.Vec2Within(mgl64.Vec2{e.Position()[0], e.Position()[2]}) {
			claim.h.HandleAttackEntity(ctx, e, force, height)
			return
		}
	}
	Wilderness.h.HandleAttackEntity(ctx, e, force, height)
}

func actuallyMovedXZ(old, new mgl64.Vec3) bool {
	return old.X() != new.X() || old.Z() != new.Z()
}

func (c *ClaimHandler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if actuallyMovedXZ(c.p.Position(), newPos) {
		c.p.SendTip(math.Round(newPos[0]), math.Round(newPos[2]))
		for _, claim := range claims {
			if !claim.area.Vec2Within(mgl64.Vec2{newPos[0], newPos[2]}) {
				claim.LeaveClaim(ctx, c.p)
			} else {
				claim.EnterClaim(ctx, c.p)
			}
		}
	}

}
