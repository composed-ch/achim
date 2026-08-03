package labels

import v3 "github.com/exoscale/egoscale/v3"

type Filterable interface {
	Name() string
	Labels() map[string]string
}

type FilterableInstance struct {
	Instance v3.Instance
}

func (n FilterableInstance) Name() string {
	return n.Instance.Name
}

func (n FilterableInstance) Labels() map[string]string {
	return n.Instance.Labels
}

func ToFilterableInstances(instances []v3.Instance) []FilterableInstance {
	result := make([]FilterableInstance, len(instances))
	for j, i := range instances {
		result[j] = FilterableInstance{Instance: i}
	}
	return result
}

func UnwrapInstances(filterables []FilterableInstance) []v3.Instance {
	instances := make([]v3.Instance, len(filterables))
	for i, f := range filterables {
		instances[i] = f.Instance
	}
	return instances
}

type FilterableNetwork struct {
	Network v3.PrivateNetwork
}

func (n FilterableNetwork) Name() string {
	return n.Network.Name
}

func (n FilterableNetwork) Labels() map[string]string {
	return n.Network.Labels
}

func ToFilterableNetworks(networks []v3.PrivateNetwork) []FilterableNetwork {
	result := make([]FilterableNetwork, len(networks))
	for i, n := range networks {
		result[i] = FilterableNetwork{Network: n}
	}
	return result
}

func UnwrapNetworks(filterables []FilterableNetwork) []v3.PrivateNetwork {
	networks := make([]v3.PrivateNetwork, len(filterables))
	for i, f := range filterables {
		networks[i] = f.Network
	}
	return networks
}
