package achim

import (
	"context"
	"fmt"
	"net"

	"github.com/composed-ch/achim/labels"
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

func ListNetworks(ctx context.Context, by string) error {
	exo := ctx.Value("exo").(*v3.Client)
	result, err := exo.ListPrivateNetworks(ctx)
	if err != nil {
		return fmt.Errorf("list networks: %w", err)
	}
	matching, err := labels.Filter(labels.ToFilterableNetworks(result.PrivateNetworks), by)
	if err != nil {
		return fmt.Errorf(`filter networks by "%s": %w`, by, err)
	}
	for _, item := range matching {
		n := item.Network
		fmt.Printf("network %s [%s,%s]:%s\n", n.Name, n.StartIP, n.EndIP, n.Netmask)
	}
	return nil
}
