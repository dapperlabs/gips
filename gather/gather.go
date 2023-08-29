package gather

import (
	"context"
	"time"

	"github.com/darron/gips/config"
	"github.com/darron/gips/core"
	"github.com/darron/gnatips-go"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"
)

var (
	loopDelay = time.Minute * 15
)

func Start(config *config.App, projects config.Projects) {
	// Start gathering data from GCP - loop forever.
	for {
		start := time.Now()
		config.Logger.Info("gather.Start")
		err := gather(config, projects)
		if err != nil {
			config.Logger.Error(err.Error())
		}
		since := time.Since(start)
		config.Logger.Info("gather.Start", "sleeping", loopDelay.String(), "loopTime", since.String())
		time.Sleep(loopDelay)
	}
}

func gather(config *config.App, projects config.Projects) error {
	// Loop through all the projects and grab all of the IPs in the various regions.
	for _, project := range projects.Projects {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		p := core.Project{
			Name: project.Name,
		}
		// Get CloudSQL IPs - these aren't organized by region.
		sqlIPs, err := getSQLIPs(ctx, config, project.Name)
		if err != nil {
			config.Logger.Error(err.Error())
		}
		if len(sqlIPs) > 0 {
			p.Regions = append(p.Regions, core.ProjectRegionIPs{
				Region: "sql",
				IPs:    sqlIPs,
			})
		}
		for _, region := range project.Regions {
			// Get NAT IPs.
			ips, err := gnatips.Get(ctx, project.Name, region)
			if err != nil {
				config.Logger.Error(err.Error())
			}
			// Let's not bother saving an empty region
			if len(ips) > 0 {
				p.Regions = append(p.Regions, core.ProjectRegionIPs{
					Region: region,
					IPs:    ips,
				})
			}
		}
		// Let's show the whole project.
		config.Logger.Info("gather", "project", p)
		// Store the project in memory.
		_, err = config.ProjectStore.Store(&p)
		if err != nil {
			config.Logger.Error(err.Error())
		}
	}
	return nil
}

func getSQLIPs(ctx context.Context, config *config.App, project string) ([]string, error) {
	var ips []string
	var err error

	// Create the Google Cloud SQL service.
	service, err := sqladmin.NewService(ctx, option.WithScopes(sqladmin.SqlserviceAdminScope))
	if err != nil {
		return ips, err
	}

	// List instances for the project ID.
	instances, err := service.Instances.List(project).Do()
	if err != nil {
		return ips, err
	}
	for _, instance := range instances.Items {
		for _, ip := range instance.IpAddresses {
			ips = append(ips, ip.IpAddress)
		}
	}
	return ips, err
}
