package selector

import (
	"time"
	"math/rand"
)

type randomBalancer struct {

}

func (r *randomBalancer) Balance(serviceName string, nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	num := rand.Intn(len(nodes))
	return nodes[num]
}

func newRandomBalancer() *randomBalancer {
	return &randomBalancer{}
}
