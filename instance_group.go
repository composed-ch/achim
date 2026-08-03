package achim

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"

	"github.com/composed-ch/achim/labels"
	v3 "github.com/exoscale/egoscale/v3"
	"gopkg.in/yaml.v3"
)

type NewInstanceGroupParam struct {
	Names     map[string]User
	Key       string
	Autostart bool
	Image     string
	Size      string
	Labels    string
	CloudInit string
}

func (p *NewInstanceGroupParam) Compile(ctx context.Context) (*InstanceGroup, error) {
	labels, err := labels.ParseLabels(p.Labels)
	if err != nil {
		return nil, fmt.Errorf("parse labels: %w", err)
	}
	existingInstances, err := ListInstances(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("list instances: %w", err)
	}
	existingNames := make(map[string]struct{}, 0)
	for _, i := range existingInstances {
		existingNames[i.Name] = struct{}{}
	}
	missingNames := make([]string, 0)
	for required := range p.Names {
		if _, ok := existingNames[required]; !ok {
			missingNames = append(missingNames, required)
		}
	}
	template, err := GetTemplateByName(ctx, p.Image)
	if err != nil {
		return nil, fmt.Errorf("find image: %w", err)
	}
	diskSizeGb := template.Size / (1024 * 1024 * 1024)
	if diskSizeGb == 0 {
		diskSizeGb = 10
	}
	instanceType, err := GetInstanceTypeBySize(ctx, p.Size)
	if err != nil {
		return nil, fmt.Errorf("find instance type: %w", err)
	}
	sshKey, err := GetSSHKeyByName(ctx, p.Key)
	if err != nil {
		return nil, fmt.Errorf("find SSH key: %w", err)
	}
	cloudInitData := make(map[string]string)
	if p.CloudInit != "" {
		cloudInitTmpl, err := parseCloudInitTemplate(p.CloudInit)
		if err != nil {
			return nil, fmt.Errorf("parse cloud init file: %w", err)
		}
		for _, name := range missingNames {
			if user, ok := p.Names[name]; ok {
				buf := bytes.NewBufferString("")
				err := cloudInitTmpl.Execute(buf, user)
				if err != nil {
					return nil, fmt.Errorf("apply cloud init template to %v: %w", user, err)
				}
				cloudInitData[name] = buf.String()
			}
		}
	}
	return &InstanceGroup{
		Names:         missingNames,
		Key:           sshKey,
		Autostart:     p.Autostart,
		Template:      template,
		InstanceType:  instanceType,
		DiskSizeGB:    diskSizeGb,
		Labels:        labels,
		CloudInitData: cloudInitData,
	}, nil
}

type InstanceGroup struct {
	Names         []string
	Key           *v3.SSHKey
	Autostart     bool
	Template      *v3.Template
	InstanceType  *v3.InstanceType
	DiskSizeGB    int64
	Labels        []labels.Label
	CloudInitData map[string]string
}

func (p *InstanceGroup) Create(ctx context.Context) error {
	exo := ctx.Value("exo").(*v3.Client)
	for _, name := range p.Names {
		genericLabels := map[string]string{"name": name}
		userData, ok := p.CloudInitData[name]
		if len(p.CloudInitData) > 0 && !ok {
			return fmt.Errorf("missing cloud init data for user")
		}
		userDataBase64 := base64.StdEncoding.EncodeToString([]byte(userData))
		_, err := exo.CreateInstance(ctx, v3.CreateInstanceRequest{
			AutoStart:    &p.Autostart,
			DiskSize:     p.DiskSizeGB,
			InstanceType: p.InstanceType,
			Labels:       labels.MergeMaps(labels.AsMap(p.Labels), genericLabels),
			Name:         name,
			SSHKey:       p.Key,
			Template:     p.Template,
			UserData:     userDataBase64,
		})
		if err != nil {
			return fmt.Errorf("create instance: %w", err)
		}
	}
	return nil
}

func parseCloudInitTemplate(yamlPath string) (*template.Template, error) {
	raw, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf(`read YAML from file "%s": %w`, yamlPath, err)
	}
	obj := make(map[string]any)
	if err := yaml.Unmarshal(raw, obj); err != nil {
		return nil, fmt.Errorf(`parse YAML content: %w`, err)
	}
	tmpl, err := template.New("cloud-init").Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf(`parse YAML as template: %w`, err)
	}
	return tmpl, nil
}
