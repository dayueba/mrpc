package selector

type Balancer interface {
	Balance(string, []*Node) *Node
}

var balancerMap = make(map[string]Balancer, 0)


const (
	Random = "random"
)

func init() {
	RegisterBalancer(Random, DefaultBalancer)
}

var DefaultBalancer = newRandomBalancer()

func RegisterBalancer(name string, balancer Balancer) {
	if balancerMap == nil {
		balancerMap = make(map[string]Balancer)
	}
	balancerMap[name] = balancer
}

func GetBalancer(name string) Balancer {
	if balancer, ok := balancerMap[name]; ok {
		return balancer
	}
	return DefaultBalancer
}
