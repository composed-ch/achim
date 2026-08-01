package achim

import (
	"context"
	"fmt"

	v3 "github.com/exoscale/egoscale/v3"
)

func FlushRecords(ctx context.Context, domain string) error {
	exo := ctx.Value("exo").(*v3.Client)
	domains, err := exo.ListDNSDomains(ctx)
	if err != nil {
		return fmt.Errorf(`list DNS domains: %w`, err)
	}
	matchingDomain, err := domains.FindDNSDomain(domain)
	if err != nil {
		return fmt.Errorf(`domain "%s" not found: %w`, domain, err)
	}
	res, err := exo.ListDNSDomainRecords(ctx, matchingDomain.ID)
	if err != nil {
		return fmt.Errorf(`list records for domain "%s": %w`, err)
	}
	for _, record := range res.DNSDomainRecords {
		if !*record.SystemRecord {
			_, err := exo.DeleteDNSDomainRecord(ctx, matchingDomain.ID, record.ID)
			if err != nil {
				return fmt.Errorf(`delete record "%s" for domain "%s": %w`, record.Name, domain, err)
			}
		}
	}
	return nil
}
