package achim

import (
	"context"
	"fmt"
	"net"
	"slices"

	"github.com/composed-ch/achim/labels"
	v3 "github.com/exoscale/egoscale/v3"
)

type NewScenarioParams struct {
	ScenarioFile string
	GroupFile    string
	KeyName      string
	Autostart    bool
	Labels       string
}

type Scenario struct {
	Name      string             `yaml:"name"`
	Instances []ScenarioInstance `yaml:"instances"`
	Networks  []ScenarioNetwork  `yaml:"networks"`
}

const (
	MinDiskSizeGB = 10
	MaxDiskSizeGB = 250
)

type ScenarioInstance struct {
	Name       string `yaml:"name"`
	Image      string `yaml:"image"`
	Size       string `yaml:"size"`
	DiskSizeGB uint   `yaml:"disk-gb"`
}

type ScenarioNetwork struct {
	Name     string `yaml:"name"`
	StartIP  string `yaml:"start-ip"`
	EndIP    string `yaml:"end-ip"`
	Netmask  string `yaml:"netmask"`
	Connects []ScenarioNetworkConnection
}

type ScenarioNetworkConnection struct {
	Instance string `yaml:"instance"`
	IP       string `yaml:"ip"`
}

func ParseScenarioFile(path string) (*Scenario, error) {
	return ParseYAMLFile[Scenario](path)
}

// ValidateInternally validates all the IP addresses in the network section and
// whether or not the connections defined refer existing instances.
func (s *Scenario) ValidateInternally() error {
	instances := make(map[string]ScenarioInstance, len(s.Instances))
	for _, instance := range s.Instances {
		instances[instance.Name] = instance
		if instance.DiskSizeGB < MinDiskSizeGB || instance.DiskSizeGB > MaxDiskSizeGB {
			return fmt.Errorf("disk size for instance %s is %d, must be within [%d;%d]",
				instance.Name, instance.DiskSizeGB, MinDiskSizeGB, MaxDiskSizeGB)
		}
	}
	for _, network := range s.Networks {
		if net.ParseIP(network.StartIP) == nil {
			return fmt.Errorf("invalid start IP %s for network %s", network.StartIP, network.Name)
		}
		if net.ParseIP(network.EndIP) == nil {
			return fmt.Errorf("invalid end IP %s for network %s", network.EndIP, network.Name)
		}
		if net.ParseIP(network.Netmask) == nil {
			return fmt.Errorf("invalid netmask %s for network %s", network.Netmask, network.Name)
		}
		for _, connects := range network.Connects {
			if net.ParseIP(connects.IP) == nil {
				return fmt.Errorf("invalid IP %s for network %s connection to %s",
					connects.IP, network.Name, connects.Instance)
			}
			if _, ok := instances[connects.Instance]; !ok {
				return fmt.Errorf("no such instance %s for network %s", connects.Instance, network.Name)
			}
		}
	}
	return nil
}

// ValidateExternally validates all the image and size definitions in the instances section.
func (s *Scenario) ValidateExternally(ctx context.Context) error {
	allowedSizes, err := GetAllowedSizes(ctx)
	if err != nil {
		return fmt.Errorf("look up allowed sizes: %w", err)
	}
	for _, instance := range s.Instances {
		_, err := GetTemplateByName(ctx, instance.Image)
		if err != nil {
			return fmt.Errorf("get image for instance %s by name %s: %w", instance.Name, instance.Image, err)
		}
		if !slices.Contains(allowedSizes, instance.Size) {
			return fmt.Errorf("no such size %s (demanded by instance %s)", instance.Size, instance.Name)
		}
	}
	return nil
}

func CreateScenario(ctx context.Context, params NewScenarioParams) error {
	exo := ctx.Value("exo").(*v3.Client)
	scenario, err := ParseScenarioFile(params.ScenarioFile)
	if err != nil {
		return fmt.Errorf(`parse scenario file "%s": %w`, params.ScenarioFile, err)
	}
	if err := scenario.ValidateInternally(); err != nil {
		return fmt.Errorf("validate scenario internally: %w", err)
	}
	if err := scenario.ValidateExternally(ctx); err != nil {
		return fmt.Errorf("validate scenario externally: %w", err)
	}
	group, err := ParseGroupFile(params.GroupFile)
	if err != nil {
		return fmt.Errorf(`parse group file "%s": %w`, params.GroupFile, err)
	}
	labels, err := labels.ParseLabels(params.Labels)
	if err != nil {
		return fmt.Errorf(`parse labels "%s": %w`, params.Labels, err)
	}
	fmt.Println(exo, scenario, group, labels)
	return nil
}
