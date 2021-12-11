package claim

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dragonfly-on-steroids/area"
	"github.com/go-gl/mathgl/mgl64"
)

var claims []Claim

func Register(c Claim) {
	claims = append(claims, c)
}

func Claims() []Claim { return claims }

type Claim interface {
	World() *world.World
	Area() area.Area
	Name() string

	Enter(p *player.Player)
	Leave(p *player.Player)

	AllowBreakBlock(p *player.Player, pos cube.Pos, drops *[]item.Stack) bool
	AllowEnter(p *player.Player) bool
	AllowAttackEntity(p *player.Player, e world.Entity, force *float64, height *float64) bool
}

type NopClaim struct{}

func (NopClaim) World() *world.World    { return nil }
func (NopClaim) Area() area.Area        { return area.Area{} }
func (NopClaim) Name() string           { return "" }
func (NopClaim) Enter(p *player.Player) {}
func (NopClaim) Leave(p *player.Player) {}

func (NopClaim) AllowBreakBlock(p *player.Player, pos cube.Pos, drops *[]item.Stack) bool {
	return true
}
func (NopClaim) AllowEnter(p *player.Player) bool                                        { return true }
func (NopClaim) AllowAttackEntity(*player.Player, world.Entity, *float64, *float64) bool { return true }

func VecInClaimXZ(vec mgl64.Vec3) (Claim, bool) {
	for _, claim := range Claims() {
		if claim.Area().Vec2Within(mgl64.Vec2{vec[0], vec[2]}) {
			return claim, true
		}
	}
	return nil, false
}
