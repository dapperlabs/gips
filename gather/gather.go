package gather

import (
	"context"
	"fmt"
	"time"

	"github.com/darron/gips/config"
	"github.com/darron/gips/core"
	"github.com/darron/gnatips-go"
)

var (
	loopDelay = time.Minute * 5
)

func Start(config *config.App, projects config.Projects) {
	// Start gathering data from GCP - loop forever.
	for {
		config.Logger.Info("gather.Start")
		err := gather(config, projects)
		if err != nil {
			config.Logger.Error(err.Error())
		}
		// Let's see all the data we have.
		data, err := config.ProjectStore.GetAll()
		if err != nil {
			config.Logger.Error(err.Error())
		}
		fmt.Printf("data: %+v\n", data)
		config.Logger.Info("gather.Start", "sleeping", loopDelay.String())
		time.Sleep(loopDelay)
	}
}

func gather(config *config.App, projects config.Projects) error {
	// Loop through all the projects and grab all of the IPs in the various regions.
	for _, project := range projects.Projects {
		p := core.Project{
			Name: project.Name,
		}
		for _, region := range project.Regions {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// TODO: Get CloudSQL IPs.
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
			config.Logger.Info("gather", "project", project.Name, "region", region, "ips", ips)
		}
		// Store the project in memory.
		_, err := config.ProjectStore.Store(&p)
		if err != nil {
			config.Logger.Error(err.Error())
		}
	}
	return nil
}
