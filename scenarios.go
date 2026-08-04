package achim

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/composed-ch/achim/labels"
	"github.com/composed-ch/goset"
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
	DiskSizeGB int64  `yaml:"disk-gb"`
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

type InstanceCreationCache struct {
	TemplatesByImageName map[string]*v3.Template
	InstanceTypesBySize  map[string]*v3.InstanceType
}

// ValidateExternally validates all the image and size definitions in the instances section.
func (s *Scenario) ValidateExternally(ctx context.Context) (*InstanceCreationCache, error) {
	cache := InstanceCreationCache{
		TemplatesByImageName: make(map[string]*v3.Template),
		InstanceTypesBySize:  make(map[string]*v3.InstanceType),
	}
	allowedSizes, err := GetAllowedSizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("look up allowed sizes: %w", err)
	}
	for _, instance := range s.Instances {
		template, err := GetTemplateByName(ctx, instance.Image)
		if err != nil {
			return nil, fmt.Errorf("get image for instance %s by name %s: %w", instance.Name, instance.Image, err)
		} else {
			cache.TemplatesByImageName[instance.Image] = template
		}
		if instanceType, ok := allowedSizes[instance.Size]; !ok {
			return nil, fmt.Errorf("no such size %s (demanded by instance %s)", instance.Size, instance.Name)
		} else {
			cache.InstanceTypesBySize[instance.Size] = instanceType
		}
	}
	return &cache, nil
}

type ScenarioSetup struct {
	Instances      map[string]ScenarioInstance
	InstanceOwners map[string]string
	Networks       map[string]ScenarioNetwork
	Attachments    map[string]map[string]net.IP
}

func (s *Scenario) CompileFor(ctx context.Context, g *Group) *ScenarioSetup {
	nInstances := len(s.Instances)
	nNetworks := len(s.Networks)
	nUsers := len(g.Users)
	setup := ScenarioSetup{
		Instances:      make(map[string]ScenarioInstance, nInstances*nUsers),
		InstanceOwners: make(map[string]string, nInstances*nUsers),
		Networks:       make(map[string]ScenarioNetwork, nNetworks*nUsers),
		Attachments:    make(map[string]map[string]net.IP),
	}
	for _, user := range g.Users {
		for _, instance := range s.Instances {
			instanceName := fmt.Sprintf("%s_%s", user.Name, instance.Name)
			setup.Instances[instanceName] = instance
			setup.InstanceOwners[instanceName] = user.Name
		}
		for _, network := range s.Networks {
			networkName := fmt.Sprintf("%s_%s", user.Name, network.Name)
			setup.Networks[networkName] = network
			for _, connect := range network.Connects {
				instanceName := fmt.Sprintf("%s_%s", user.Name, connect.Instance)
				if _, ok := setup.Attachments[instanceName]; !ok {
					setup.Attachments[instanceName] = make(map[string]net.IP)
				}
				setup.Attachments[instanceName][networkName] = net.ParseIP(connect.IP)
			}
		}
	}
	return &setup
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
	cache, err := scenario.ValidateExternally(ctx)
	if err != nil {
		return fmt.Errorf("validate scenario externally: %w", err)
	}
	group, err := ParseGroupFile(params.GroupFile)
	if err != nil {
		return fmt.Errorf(`parse group file "%s": %w`, params.GroupFile, err)
	}
	key, err := GetSSHKeyByName(ctx, params.KeyName)
	if err != nil {
		return fmt.Errorf(`get SSH key by name "%s": %w`, params.KeyName, err)
	}
	sharedLabels, err := labels.ParseLabels(params.Labels)
	if err != nil {
		return fmt.Errorf(`parse labels "%s": %w`, params.Labels, err)
	}
	setup := scenario.CompileFor(ctx, group)
	if err := uniqueInstances(ctx, setup); err != nil {
		return err
	}
	if err := uniqueNetworks(ctx, setup); err != nil {
		return err
	}
	instanceIDsByName := make(map[string]v3.UUID)
	for instanceName, details := range setup.Instances {
		owner := setup.InstanceOwners[instanceName]
		specificLabels := map[string]string{
			"owner":    owner,
			"scenario": scenario.Name,
			"group":    group.Name,
		}
		instanceLabels := labels.MergeMaps(labels.AsMap(sharedLabels), specificLabels)
		instanceType := cache.InstanceTypesBySize[details.Size]
		template := cache.TemplatesByImageName[details.Image]
		req := v3.CreateInstanceRequest{
			AutoStart:    &params.Autostart,
			DiskSize:     details.DiskSizeGB,
			InstanceType: instanceType,
			Labels:       instanceLabels,
			Name:         instanceName,
			SSHKey:       key,
			Template:     template,
		}
		op, err := exo.CreateInstance(ctx, req)
		if err != nil {
			return fmt.Errorf(`create instance "%s": %w`, instanceName, err)
		}
		instanceIDsByName[instanceName] = op.Reference.ID
		fmt.Printf("created instance %s for %s\n", instanceName, owner)
	}
	networkIDsByName := make(map[string]v3.UUID)
	for networkName, details := range setup.Networks {
		req := v3.CreatePrivateNetworkRequest{
			Name:        networkName,
			Description: "", // TODO: incorporate network description into YAML file format
			StartIP:     net.ParseIP(details.StartIP),
			EndIP:       net.ParseIP(details.EndIP),
			Netmask:     net.ParseIP(details.Netmask),
			// TODO: labels as above (requires additional NetworkOwners map on setup)
		}
		op, err := exo.CreatePrivateNetwork(ctx, req)
		if err != nil {
			return fmt.Errorf(`create network "%s": %w`, networkName, err)
		}
		networkIDsByName[networkName] = op.Reference.ID
		fmt.Printf("created network %s\n", networkName /* TODO: output owner */)
	}
	for instanceName, attachments := range setup.Attachments {
		for networkName, ip := range attachments {
			networkID := networkIDsByName[networkName]
			instanceID := instanceIDsByName[instanceName]
			req := v3.AttachInstanceToPrivateNetworkRequest{
				IP: ip,
				Instance: &v3.AttachInstanceToPrivateNetworkRequestInstance{
					ID: instanceID,
				},
			}
			_, err := exo.AttachInstanceToPrivateNetwork(ctx, networkID, req)
			if err != nil {
				return fmt.Errorf("attach network %s to instance %s with IP %s: %w", networkName, instanceName, ip)
			}
			fmt.Printf("attached network %s to instance %s with IP %s\n", networkName, instanceName, ip)
		}
	}
	return nil
}

func uniqueInstances(ctx context.Context, setup *ScenarioSetup) error {
	existingInstances, err := ListInstances(ctx, "")
	if err != nil {
		return fmt.Errorf("list existing instances: %w", err)
	}
	existingInstanceNames := make([]string, len(existingInstances))
	for i, instance := range existingInstances {
		existingInstanceNames[i] = instance.Name
	}
	scenarioInstanceNames := make([]string, 0)
	for name := range setup.Instances {
		scenarioInstanceNames = append(scenarioInstanceNames, name)
	}
	conflictingInstanceNames := goset.From(existingInstanceNames).Inter(goset.From(scenarioInstanceNames))
	if len(conflictingInstanceNames.Entries) > 0 {
		return fmt.Errorf("instances already exist: %s", strings.Join(conflictingInstanceNames.Slice(), ", "))
	}
	return nil
}

func uniqueNetworks(ctx context.Context, setup *ScenarioSetup) error {
	existingNetworks, err := GetNetworks(ctx, "")
	if err != nil {
		return fmt.Errorf("list existing instances: %w", err)
	}
	existingNetworkNames := make([]string, len(existingNetworks))
	for i, network := range existingNetworks {
		existingNetworkNames[i] = network.Name
	}
	scenarioNetworkNames := make([]string, 0)
	for name := range setup.Instances {
		scenarioNetworkNames = append(scenarioNetworkNames, name)
	}
	conflictingNetworkNames := goset.From(existingNetworkNames).Inter(goset.From(scenarioNetworkNames))
	if len(conflictingNetworkNames.Entries) > 0 {
		return fmt.Errorf("network already exist: %s", strings.Join(conflictingNetworkNames.Slice(), ", "))
	}
	return nil
}
