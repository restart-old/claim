package main

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dragonfly-on-steroids/area"
	"github.com/dragonfly-on-steroids/claim"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sirupsen/logrus"
)

type testClaim struct {
	w *world.World
	claim.NopClaim
}

func (testClaim) Area() area.Area {
	return area.NewArea(mgl64.Vec2{0, 0}, mgl64.Vec2{10, 10})
}
func (testClaim) Name() string          { return "claim" }
func (t testClaim) World() *world.World { return t.w }
func (testClaim) AllowBreakBlock(p *player.Player, pos cube.Pos, drops *[]item.Stack) bool {
	return false
}

func (w testClaim) Enter(p *player.Player) {
	p.Messagef("§eNow entering: §ctest §e(§cDeathban§e)")
}
func (w testClaim) Leave(p *player.Player) {
	p.Messagef("§eNow leaving: §ctest §e(§cDeathban§e)")
}

func main() {
	c := server.DefaultConfig()
	c.Players.SaveData = false
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel
	s := server.New(&c, log)
	s.Start()

	claim.Register(testClaim{w: s.World()})
	wild := claim.NewWilderness(s.World(), "§eNow entering: §7The Wilderness §e(§cDeathban§e)", "§eNow leaving: §7The Wilderness §e(§cDeathban§e)")
	for {
		p, err := s.Accept()
		if err != nil {
			return
		}
		p.Handle(claim.NewClaimHandler(p, wild))
	}
}
