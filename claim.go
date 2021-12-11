package claim

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/physics"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

var claims []Claim

func Register(c Claim) {
	claims = append(claims, c)
}

func Claims() []Claim { return claims }

type Claim interface {
	World() *world.World
	AABB() physics.AABB
	Name() string

	Enter(p *player.Player)
	Leave(p *player.Player)

	AllowBreakBlock(p *player.Player, pos cube.Pos, drops *[]item.Stack) bool
	AllowEnter(p *player.Player) bool
	AllowAttackEntity(p *player.Player, e world.Entity, force *float64, height *float64) bool
}

type NopClaim struct{}

func (NopClaim) World() *world.World    { return nil }
func (NopClaim) AABB() physics.AABB     { return physics.AABB{} }
func (NopClaim) Name() string           { return "" }
func (NopClaim) Enter(p *player.Player) {}
func (NopClaim) Leave(p *player.Player) {}

func (NopClaim) AllowBreakBlock(p *player.Player, pos cube.Pos, drops *[]item.Stack) bool {
	return true
}
func (NopClaim) AllowEnter(p *player.Player) bool                                        { return true }
func (NopClaim) AllowAttackEntity(*player.Player, world.Entity, *float64, *float64) bool { return true }

func inOrEqual(vec mgl64.Vec3, aabb physics.AABB) bool {
	if vec[0] < aabb.Min()[0] || vec[0] > aabb.Max()[0] {
		return false
	}
	return vec[2] >= aabb.Min()[2] && vec[2] <= aabb.Max()[2]
}
func PosInClaim(pos mgl64.Vec3) (Claim, bool) {
	for _, claim := range Claims() {
		if inOrEqual(pos, claim.AABB()) {
			return claim, true
		}
	}
	return nil, false
}
