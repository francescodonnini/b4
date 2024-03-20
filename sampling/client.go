package sampling

import "b4/shared"

type Client interface {
	Exchange(view PView, dest shared.Node) (PView, error)
}
