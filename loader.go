package claim

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Loader interface {
	LoadWithPos(mgl64.Vec3) (*Claim, error)
}
