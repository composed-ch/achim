package achim

import (
	"context"
	"fmt"

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

type ScenarioInstance struct {
}

type ScenarioNetwork struct {
}

func ParseScenarioFile(path string) (*Scenario, error) {
	return nil, nil
}

func CreateScenario(ctx context.Context, params NewScenarioParams) error {
	exo := ctx.Value("exo").(*v3.Client)
	scenario, err := ParseScenarioFile(params.ScenarioFile)
	if err != nil {
		return fmt.Errorf(`parse scenario file "%s": %w`, params.ScenarioFile, err)
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
