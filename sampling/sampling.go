package sampling

import (
	"b4/shared"
)

type PeerSampling interface {
	GetRandom() (shared.Node, bool)
}
