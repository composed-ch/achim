package achim

import (
	"context"
	"fmt"

	"github.com/composed-ch/goset"
	v3 "github.com/exoscale/egoscale/v3"
)

func FlushRecords(ctx context.Context, domain string) error {
	exo := ctx.Value("exo").(*v3.Client)
	matchingDomain, err := getDNSDomain(ctx, domain)
	if err != nil {
		return fmt.Errorf("get DNS domain: %w", err)
	}
	res, err := exo.ListDNSDomainRecords(ctx, matchingDomain.ID)
	if err != nil {
		return fmt.Errorf(`list records for domain "%s": %w`, domain, err)
	}
	for _, record := range res.DNSDomainRecords {
		if !*record.SystemRecord {
			_, err := exo.DeleteDNSDomainRecord(ctx, matchingDomain.ID, record.ID)
			if err != nil {
				return fmt.Errorf(`delete record "%s" for domain "%s": %w`, record.Name, domain, err)
			} else {
				fmt.Printf("deleted DNS record %s.%s=%s\n", record.Name, domain, record.Content)
			}
		}
	}
	return nil
}

type DNSEntry struct {
	Name    string
	Content string
}

func SyncRecords(ctx context.Context, domain string) error {
	exo := ctx.Value("exo").(*v3.Client)
	matchingDomain, err := getDNSDomain(ctx, domain)
	if err != nil {
		return fmt.Errorf("get DNS domain: %w", err)
	}
	res, err := exo.ListDNSDomainRecords(ctx, matchingDomain.ID)
	if err != nil {
		return fmt.Errorf(`list records for domain "%s": %w`, domain, err)
	}
	existing := make([]DNSEntry, 0)
	existingId := make(map[DNSEntry]v3.UUID)
	for _, record := range res.DNSDomainRecords {
		if *record.SystemRecord {
			continue
		}
		entry := DNSEntry{Name: record.Name, Content: record.Content}
		existing = append(existing, entry)
		existingId[entry] = record.ID
	}
	instances, err := exo.ListInstances(ctx)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	required := make([]DNSEntry, 0)
	for _, instance := range instances.Instances {
		required = append(required, DNSEntry{Name: instance.Name, Content: instance.PublicIP.String()})
	}
	existingSet := goset.From(existing)
	requiredSet := goset.From(required)
	toBeCreated := requiredSet.Diff(existingSet)
	toBeDeleted := existingSet.Diff(requiredSet)
	for _, createEntry := range toBeCreated.Slice() {
		var ttl int64 = 3600
		_, err := exo.CreateDNSDomainRecord(ctx, matchingDomain.ID, v3.CreateDNSDomainRecordRequest{
			Content: createEntry.Content,
			Name:    createEntry.Name,
			Ttl:     ttl,
			Type:    v3.CreateDNSDomainRecordRequestTypeA,
		})
		if err != nil {
			return fmt.Errorf(`create DNS record %v for domain "%s": %w`, createEntry, domain, err)
		} else {
			fmt.Printf("created DNS %s record %s.%s=%s TTL=%d\n",
				v3.CreateDNSDomainRecordRequestTypeA, createEntry.Name, domain, createEntry.Content, ttl)
		}
	}
	for _, deleteEntry := range toBeDeleted.Slice() {
		id, ok := existingId[deleteEntry]
		if !ok {
			return fmt.Errorf("unknown ID for record %v", deleteEntry)
		}
		_, err := exo.DeleteDNSDomainRecord(ctx, matchingDomain.ID, id)
		if err != nil {
			return fmt.Errorf(`delete DNS record %v (id: %v) for domain "%s": %w`, deleteEntry, id, domain, err)
		} else {
			fmt.Printf("deleted DNS record %s.%s=%s\n", deleteEntry.Name, domain, deleteEntry.Content)
		}
	}
	return nil
}

func getDNSDomain(ctx context.Context, name string) (*v3.DNSDomain, error) {
	exo := ctx.Value("exo").(*v3.Client)
	domains, err := exo.ListDNSDomains(ctx)
	if err != nil {
		return nil, fmt.Errorf(`list DNS domains: %w`, err)
	}
	domain, err := domains.FindDNSDomain(name)
	if err != nil {
		return nil, fmt.Errorf(`domain "%s" not found: %w`, name, err)
	}
	return &domain, nil
}
