package main

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/dragonfly-on-steroids/area"
	"github.com/dragonfly-on-steroids/claim"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sirupsen/logrus"
)

func main() {
	c := server.DefaultConfig()
	c.Players.SaveData = false
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel
	s := server.New(&c, log)
	s.Start()

	cl := claim.NewClaim("test", s.World(), area.NewArea(mgl64.Vec2{10, 10}, mgl64.Vec2{20, 20}))
	cl.Handle(&ClaimHandler{c: cl})
	claim.Store(cl)
	claim.Wilderness.Handle(&ClaimHandler{c: claim.Wilderness})
	for {
		p, err := s.Accept()
		if err != nil {
			return
		}
		p.Handle(claim.NewClaimHandler(p))
	}
}

type ClaimHandler struct {
	claim.NopHandler
	c *claim.Claim
}

func (c *ClaimHandler) HandleEnterClaim(ctx *event.Context, p *player.Player) {
	p.Message("entered", c.c.Name())
}
func (c *ClaimHandler) HandleLeaveClaim(ctx *event.Context, p *player.Player) {
	p.Message("left", c.c.Name())
}
func (c *ClaimHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	if c.c != claim.Wilderness {
		ctx.Cancel()
	}
}
