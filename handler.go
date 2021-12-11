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

// ClaimHandler is the handler which is used to handle:
// When a block is broken in a claim.
// When a player enters or leaves a claim.
// And when a player is hurt in a claim.
type ClaimHandler struct {
	player.NopHandler
	p    *player.Player
	wild wilderness
}

// NewClaimHandler returns a new *ClaimHandler.
func NewClaimHandler(p *player.Player, wild wilderness) *ClaimHandler {
	return &ClaimHandler{
		p:    p,
		wild: wild,
	}
}

// HandleBlockBreak handles when a block is broken,
// and cancels the event if breaking blocks are not allowed in the claim they were broken in.
func (c *ClaimHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	if claim, ok := VecInClaimXZ(pos.Vec3()); ok {
		if !claim.AllowBreakBlock(c.p, pos, drops) {
			ctx.Cancel()
		}
	}
}

func canHurt(p *player.Player, e world.Entity, force, height *float64) bool {
	if claim, ok := LoadPlayerClaim(p); ok {
		return claim.AllowAttackEntity(p, e, force, height)
	}
	return true
}

func (c *ClaimHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64) {
	if !canHurt(c.p, e, force, height) {
		ctx.Cancel()
	}
}

func actuallyMovedXZ(old, new mgl64.Vec3) bool {
	oldX := math.Round(old.X())
	oldZ := math.Round(old.Z())
	newX := math.Round(new.X())
	newZ := math.Round(new.Z())
	return oldX != newX || oldZ != newZ
}

func oldAndNewClaim(p *player.Player, newPos mgl64.Vec3) (Claim, Claim) {
	oldClaim, _ := LoadPlayerClaim(p)
	newClaim, _ := VecInClaimXZ(newPos)
	return oldClaim, newClaim
}

func finalClaims(wild wilderness, old, new Claim) (Claim, Claim) {
	if old == nil {
		old = wild
	}
	if new == nil {
		new = wild
	}
	return old, new
}

func (c *ClaimHandler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if actuallyMovedXZ(c.p.Position(), newPos) {
		c.p.SendTip(math.Round(newPos[0]), math.Round(newPos[1]), math.Round(newPos[2]))
		old, new := oldAndNewClaim(c.p, newPos)
		old, new = finalClaims(c.wild, old, new)
		if old != new {
			if new.AllowEnter(c.p) {
				old.Leave(c.p)
				new.Enter(c.p)
				StorePlayerClaim(c.p, new)
			} else {
				ctx.Cancel()
				c.p.Teleport(c.p.Position())
			}
		}

	}

}
