package achim

import (
	"context"
	"fmt"
	"html/template"
	"maps"
	"net/http"
	"os"
	"slices"

	"github.com/composed-ch/achim/labels"
	"github.com/composed-ch/achim/templates"
	v3 "github.com/exoscale/egoscale/v3"
)

type NewInstanceParams struct {
	Name      string
	Key       string
	Autostart bool
	Image     string
	Size      string
	Labels    string
	CloudInit string
}

func CreateInstance(ctx context.Context, params NewInstanceParams) error {
	singleUserData := User{Name: params.Name, Email: "", SSHKey: ""}
	newInstancesParam := NewInstanceGroupParam{
		Names:     map[string]User{params.Name: singleUserData},
		Key:       params.Key,
		Autostart: params.Autostart,
		Image:     params.Image,
		Size:      params.Size,
		Labels:    params.Labels,
		CloudInit: params.CloudInit,
	}
	newInstanceGroup, err := newInstancesParam.Compile(ctx)
	if err != nil {
		return fmt.Errorf("compile instance group %v: %w", params, err)
	}
	return newInstanceGroup.Create(ctx)
}

func DestroyInstances(ctx context.Context, by string) error {
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	exo := ctx.Value("exo").(*v3.Client)
	for _, instance := range instances {
		_, err := exo.DeleteInstance(ctx, instance.ID)
		if err != nil {
			return fmt.Errorf(`destroy instance "%s": %w`, instance.ID, err)
		} else {
			fmt.Printf("destroyed instance %s\n", instance.Name)
		}
	}
	return nil
}

func ExportInventory(ctx context.Context, file, by string) error {
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf(`create file "%s": %w`, file, err)
	} else {
		defer f.Close()
	}
	inventory := make(map[string]map[string][]string)
	inventory["group"] = make(map[string][]string)
	inventory["name"] = make(map[string][]string)
	for _, instance := range instances {
		ip := instance.PublicIP.String()
		name := instance.Name
		inventory["name"][name] = append(inventory["name"][name], ip)
		if group, ok := instance.Labels["group"]; ok {
			inventory["group"][group] = append(inventory["group"][group], ip)
		}
	}
	for _, table := range inventory {
		for value, ips := range table {
			fmt.Fprintf(f, "[%s]\n", value)
			for _, ip := range ips {
				fmt.Fprintln(f, ip)
			}
			fmt.Fprintln(f, "")
		}
	}
	fmt.Printf("inventory created under %s\n", file)
	return nil
}

func EmbiggenDisk(ctx context.Context, by string, gb int64) error {
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	exo := ctx.Value("exo").(*v3.Client)
	for _, instance := range instances {
		if instance.DiskSize > gb {
			return fmt.Errorf(`instance %s "%s" disk size is %d GB; cannot shrink disk`,
				instance.ID, instance.Name, instance.DiskSize)
		}
		if instance.State != v3.InstanceStateStopped {
			return fmt.Errorf(`instance %s "%s" must be stopped but is %s`,
				instance.ID, instance.Name, instance.State)
		}
		_, err := exo.ResizeInstanceDisk(ctx, instance.ID, v3.ResizeInstanceDiskRequest{
			DiskSize: gb,
		})
		if err != nil {
			return fmt.Errorf(`resize disk of instance %s "%s": %w`,
				instance.ID, instance.Name, err)
		} else {
			fmt.Printf("resized disk of instance %s to %dGB\n", instance.Name, gb)
		}
	}
	return nil
}

func ScaleInstances(ctx context.Context, by string, size string) error {
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	exo := ctx.Value("exo").(*v3.Client)
	types, err := ListInstanceTypes(ctx, "standard")
	if err != nil {
		return fmt.Errorf("list standard instance types: %w", err)
	}
	allowedSizes := make(map[string]*v3.InstanceType, len(types))
	for _, t := range types {
		allowedSizes[string(t.Size)] = &t
	}
	if _, ok := allowedSizes[size]; !ok {
		return fmt.Errorf(`size "%s" is not allowed, use any of %v`,
			size, slices.Collect(maps.Keys(allowedSizes)))
	}
	for _, instance := range instances {
		if instance.State != v3.InstanceStateStopped {
			return fmt.Errorf(`instance %s "%s" must be stopped but is %s`,
				instance.ID, instance.Name, instance.State)
		}
		_, err = exo.ScaleInstance(ctx, instance.ID, v3.ScaleInstanceRequest{
			InstanceType: allowedSizes[size],
		})
		if err != nil {
			return fmt.Errorf(`scale instance %s "%s" to size %s: %w`,
				instance.ID, instance.Name, size, err)
		} else {
			fmt.Printf("resized instance %s to %s\n", instance.Name, size)
		}
	}
	return nil
}

func ListInstances(ctx context.Context, by string) ([]v3.Instance, error) {
	exo := ctx.Value("exo").(*v3.Client)
	result, err := exo.ListInstances(ctx)
	if err != nil {
		return nil, fmt.Errorf("list instances: %w", err)
	}
	instances := make([]v3.Instance, 0)
	for _, item := range result.Instances {
		instance, err := getInstanceByID(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf(`get instance by ID "%s": %w`, item.ID, err)
		}
		instances = append(instances, *instance)
	}
	matching, err := labels.Filter(labels.ToFilterableInstances(instances), by)
	if err != nil {
		return nil, fmt.Errorf("filter instances: %w", err)
	}
	return labels.UnwrapInstances(matching), err
}

func ListInstanceTypes(ctx context.Context, family string) ([]v3.InstanceType, error) {
	exo := ctx.Value("exo").(*v3.Client)
	res, err := exo.ListInstanceTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list instance types: %w", err)
	}
	var result []v3.InstanceType
	for _, it := range res.InstanceTypes {
		if !*it.Authorized || string(it.Family) != family {
			continue
		}
		result = append(result, it)
	}
	return result, nil
}

func GetInstanceTypeBySize(ctx context.Context, size string) (*v3.InstanceType, error) {
	exo := ctx.Value("exo").(*v3.Client)
	res, err := exo.ListInstanceTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list instance types: %w", err)
	}
	for _, it := range res.InstanceTypes {
		if string(it.Size) == size {
			return &it, nil
		}
	}
	return nil, fmt.Errorf(`no instance type for size "%s" found`, size)
}

func LabelInstances(ctx context.Context, label, value, by string) error {
	exo := ctx.Value("exo").(*v3.Client)
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf(`list instances to label: %w`, err)
	}
	for _, instance := range instances {
		labels := make(map[string]string, 0)
		maps.Copy(labels, instance.Labels)
		labels[label] = value
		update := v3.UpdateInstanceRequest{Labels: labels}
		if _, err := exo.UpdateInstance(ctx, instance.ID, update); err != nil {
			return fmt.Errorf(`label instance %s: %w`, instance.ID, err)
		} else {
			fmt.Printf("labeled instance %s with %s=%s\n", instance.Name, label, value)
		}
	}
	return nil
}

func ProtectInstances(ctx context.Context, by string) error {
	return changeInstanceProtection(ctx, by, true)
}

func DeprotectInstances(ctx context.Context, by string) error {
	return changeInstanceProtection(ctx, by, false)
}

func StartInstances(ctx context.Context, by string) error {
	return changeInstanceState(ctx, by, true)
}

func StopInstances(ctx context.Context, by string) error {
	return changeInstanceState(ctx, by, false)
}

func FormatInstance(instance v3.Instance) string {
	return fmt.Sprintf("%-40s %-32s %15s %-10s", instance.ID, instance.Name, instance.PublicIP, instance.State)
}

func Overview(ctx context.Context, by, file string) error {
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	entries := make([]templates.OverviewEntry, len(instances))
	for i, instance := range instances {
		var owner string
		for k, v := range instance.Labels {
			if k == "owner" {
				owner = v
			}
		}
		ip := instance.PublicIP.String()
		entries[i] = templates.OverviewEntry{
			Owner:      owner,
			HostName:   instance.Name,
			IPAddress:  ip,
			SSHCommand: fmt.Sprintf("ssh %s@%s", "user", ip),
		}
	}
	data := templates.OverviewData{
		Entries:   entries,
		Selection: by,
	}
	tmpl, err := template.New("overview").Parse(templates.Overview)
	if err != nil {
		return fmt.Errorf("parse overview template: %w", err)
	}
	var html *os.File
	if file != "" {
		f, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("create %s: %w", file, err)
		}
		defer f.Close()
		html = f
	} else {
		html = os.Stdout
	}
	if err := tmpl.Execute(html, data); err != nil {
		return fmt.Errorf("execute overview template: %w", err)
	}
	fmt.Printf("overview created under %s\n", file)
	return nil
}

func ProbeInstances(ctx context.Context, by, suffix, domain string, secure bool) error {
	exo := ctx.Value("exo").(*v3.Client)
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	var reverseDNS map[string]string
	if domain != "" {
		domain, err := getDNSDomain(ctx, domain)
		if err != nil {
			return fmt.Errorf(`get domain %s: %w`, domain, err)
		}
		res, err := exo.ListDNSDomainRecords(ctx, domain.ID)
		if err != nil {
			return fmt.Errorf(`list records for domain %s: %w`, domain, err)
		}
		reverseDNS = make(map[string]string, len(res.DNSDomainRecords))
		for _, record := range res.DNSDomainRecords {
			reverseDNS[record.Content] = record.Name
		}
	}
	for _, instance := range instances {
		proto := "http"
		if secure {
			proto = "https"
		}
		var url string
		if subdomain, ok := reverseDNS[instance.PublicIP.String()]; ok {
			url = fmt.Sprintf("%s://%s.%s/%s", proto, subdomain, domain, suffix)
		} else {
			url = fmt.Sprintf("%s://%s/%s", proto, instance.PublicIP, suffix)
		}
		res, err := http.Get(url)
		if err != nil {
			fmt.Printf("%-32s GET %s\tERR\n", instance.Name, url)
			continue
		}
		fmt.Printf("%-32s GET %-50s\t%d\n", instance.Name, url, res.StatusCode)
	}
	return nil
}

func changeInstanceProtection(ctx context.Context, by string, protect bool) error {
	exo := ctx.Value("exo").(*v3.Client)
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf(`list instances to change protection to %v: %w`, protect, err)
	}
	for _, instance := range instances {
		if protect {
			if _, err := exo.AddInstanceProtection(ctx, instance.ID); err != nil {
				return fmt.Errorf(`protect instance %s: %w`, instance.ID, err)
			} else {
				fmt.Printf("protected instance %s\n", instance.Name)
			}
		} else {
			if _, err := exo.RemoveInstanceProtection(ctx, instance.ID); err != nil {
				return fmt.Errorf(`deprotect instance %s: %w`, instance.ID, err)
			} else {
				fmt.Printf("deprotected instance %s\n", instance.Name)
			}
		}
	}
	return nil
}

func changeInstanceState(ctx context.Context, by string, up bool) error {
	exo := ctx.Value("exo").(*v3.Client)
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf(`list instances to change up state to %v: %w`, up, err)
	}
	for _, instance := range instances {
		if up {
			if _, err := exo.StartInstance(ctx, instance.ID, v3.StartInstanceRequest{}); err != nil {
				return err
			} else {
				fmt.Printf("started instance %s\n", instance.Name)
			}
		} else {
			if _, err := exo.StopInstance(ctx, instance.ID); err != nil {
				return err
			} else {
				fmt.Printf("stopped instance %s\n", instance.Name)
			}
		}
	}
	return nil
}

func getInstanceByID(ctx context.Context, id v3.UUID) (*v3.Instance, error) {
	exo := ctx.Value("exo").(*v3.Client)
	instance, err := exo.GetInstance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(`get instance by ID "%v": %w`, id, err)
	}
	return instance, nil
}
