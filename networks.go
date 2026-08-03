package achim

import (
	"context"
	"fmt"
	"net"

	"github.com/composed-ch/achim/labels"
	"github.com/composed-ch/goset"
	v3 "github.com/exoscale/egoscale/v3"
)

type NewNetworkParams struct {
	Name        string
	Description string
	Labels      string
	StartIP     string
	EndIP       string
	Netmask     string
}

func (n *NewNetworkParams) ParseIPs() (*NetworkIPConfig, error) {
	startIP := net.ParseIP(n.StartIP)
	if startIP == nil {
		return nil, fmt.Errorf(`parse start IP "%s": error`, n.StartIP)
	}
	endIP := net.ParseIP(n.EndIP)
	if endIP == nil {
		return nil, fmt.Errorf(`parse end IP "%s": error`, n.EndIP)
	}
	netmask := net.ParseIP(n.Netmask)
	if netmask == nil {
		return nil, fmt.Errorf(`parse netmask "%s": error`, n.Netmask)
	}
	return &NetworkIPConfig{
		StartIP: startIP,
		EndIP:   endIP,
		Netmask: netmask,
	}, nil
}

type NetworkIPConfig struct {
	StartIP net.IP
	EndIP   net.IP
	Netmask net.IP
}

func CreateNetwork(ctx context.Context, network NewNetworkParams) error {
	exo := ctx.Value("exo").(*v3.Client)
	ipConfig, err := network.ParseIPs()
	if err != nil {
		return fmt.Errorf("parse IPs of %v: %w", network, err)
	}
	var parsedLabels []labels.Label
	if network.Labels != "" {
		parsedLabels, err = labels.ParseLabels(network.Labels)
		if err != nil {
			return fmt.Errorf(`parse labels "%s": %w`, network.Labels, err)
		}
	}
	r := v3.CreatePrivateNetworkRequest{
		Name:        network.Name,
		Description: network.Description,
		StartIP:     ipConfig.StartIP,
		EndIP:       ipConfig.EndIP,
		Netmask:     ipConfig.Netmask,
		Labels:      labels.AsMap(parsedLabels),
	}
	_, err = exo.CreatePrivateNetwork(ctx, r)
	if err != nil {
		return fmt.Errorf("create network %v: %w", network, err)
	}
	fmt.Printf("created network %s [%s,%s]:%s\n", network.Name, network.StartIP, network.EndIP, network.Netmask)
	return nil
}

func GetNetworks(ctx context.Context, by string) ([]v3.PrivateNetwork, error) {
	exo := ctx.Value("exo").(*v3.Client)
	result, err := exo.ListPrivateNetworks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list networks: %w", err)
	}
	matching, err := labels.Filter(labels.ToFilterableNetworks(result.PrivateNetworks), by)
	if err != nil {
		return nil, fmt.Errorf(`filter networks by "%s": %w`, by, err)
	}
	return labels.UnwrapNetworks(matching), nil
}

func ListNetworks(ctx context.Context, by string) error {
	networks, err := GetNetworks(ctx, by)
	if err != nil {
		return fmt.Errorf(`get networks by "%s": %w`, by, err)
	}
	for _, n := range networks {
		fmt.Printf("network %s [%s,%s]:%s\n", n.Name, n.StartIP, n.EndIP, n.Netmask)
	}
	return nil
}

func Attach(ctx context.Context, network, instance, ip string) error {
	exo := ctx.Value("exo").(*v3.Client)
	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return fmt.Errorf(`parse IP address "%s": error`, ip)
	}
	matchingInstances, err := ListInstances(ctx, fmt.Sprintf("name=%s", instance))
	if err != nil {
		return fmt.Errorf(`list instances for name="%s": %w`, instance, err)
	} else if len(matchingInstances) != 1 {
		return fmt.Errorf(`expected to get one instance for name="%s", got %d`, instance, len(matchingInstances))
	}
	networksResposne, err := exo.ListPrivateNetworks(ctx)
	if err != nil {
		return fmt.Errorf("list networks: %w", err)
	}
	matchingNetwork, err := networksResposne.FindPrivateNetwork(network)
	if err != nil {
		return fmt.Errorf(`find network by name "%s": %w`, network, err)
	}
	_, err = exo.AttachInstanceToPrivateNetwork(ctx, matchingNetwork.ID, v3.AttachInstanceToPrivateNetworkRequest{
		IP:       ipAddress,
		Instance: &v3.AttachInstanceToPrivateNetworkRequestInstance{ID: matchingInstances[0].ID},
	})
	if err != nil {
		return fmt.Errorf(`attach instance "%s" to network "%s" with IP "%s": %w`, instance, network, ip, err)
	}
	fmt.Printf("attached instance %s to network %s with IP %s\n", instance, network, ip)
	return nil
}

type NetworkPredicate func(network v3.PrivateNetwork) bool

func DeleteNetworks(ctx context.Context, predicate NetworkPredicate) ([]v3.PrivateNetwork, error) {
	exo := ctx.Value("exo").(*v3.Client)
	networkResponse, err := exo.ListPrivateNetworks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list networks: %w", err)
	}
	deleted := make([]v3.PrivateNetwork, 0)
	for _, network := range networkResponse.PrivateNetworks {
		if predicate(network) {
			_, err := exo.DeletePrivateNetwork(ctx, network.ID)
			if err != nil {
				return deleted, fmt.Errorf("delete network %v: %w", network, err)
			}
			deleted = append(deleted, network)
		}
	}
	return deleted, nil
}

func CleanupNetworks(ctx context.Context) error {
	instances, err := ListInstances(ctx, "")
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	networks, err := GetNetworks(ctx, "")
	if err != nil {
		return fmt.Errorf("list networks: %w", err)
	}
	allNetworkIds := make([]v3.UUID, len(networks))
	for i, n := range networks {
		allNetworkIds[i] = n.ID
	}
	usedNetworkIds := make([]v3.UUID, 0)
	for _, instance := range instances {
		for _, n := range instance.PrivateNetworks {
			usedNetworkIds = append(usedNetworkIds, n.ID)
		}
	}
	allNetworksSet := goset.From(allNetworkIds)
	usedNetworksSet := goset.From(usedNetworkIds)
	orphanedNetworksSet := allNetworksSet.Diff(usedNetworksSet)
	deleted, err := DeleteNetworks(ctx, func(network v3.PrivateNetwork) bool {
		_, ok := orphanedNetworksSet.Entries[network.ID]
		return ok
	})
	if len(deleted) > 0 {
		for _, d := range deleted {
			fmt.Printf("deleted orphaned network %s [%s:%s]:%s\n", d.Name, d.StartIP, d.EndIP, d.Netmask)
		}
	}
	if err != nil {
		return fmt.Errorf("cleanup orphaned networks: %w", err)
	}
	return nil
}

func DestroyNetworks(ctx context.Context, by string) error {
	matchingNetworks, err := GetNetworks(ctx, by)
	if err != nil {
		return fmt.Errorf(`filter networks by "%s": %w`, by, err)
	}
	matchingNetworkIds := make(map[v3.UUID]struct{})
	for _, network := range matchingNetworks {
		matchingNetworkIds[network.ID] = struct{}{}
	}
	deleted, err := DeleteNetworks(ctx, func(network v3.PrivateNetwork) bool {
		_, ok := matchingNetworkIds[network.ID]
		return ok
	})
	if len(deleted) > 0 {
		for _, d := range deleted {
			fmt.Printf("deleted network %s [%s:%s]:%s\n", d.Name, d.StartIP, d.EndIP, d.Netmask)
		}
	}
	if err != nil {
		return fmt.Errorf(`destroy networks by "%s": %w`, by, err)
	}
	return nil
}

func FlushNetworks(ctx context.Context) error {
	deleted, err := DeleteNetworks(ctx, func(network v3.PrivateNetwork) bool { return true })
	if len(deleted) > 0 {
		for _, d := range deleted {
			fmt.Printf("deleted network %s [%s:%s]:%s\n", d.Name, d.StartIP, d.EndIP, d.Netmask)
		}
	}
	if err != nil {
		return fmt.Errorf("flush networks: %w", err)
	}
	return nil
}
