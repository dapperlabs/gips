package gather

import (
	"context"
	"time"

	"github.com/darron/gips/config"
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
		config.Logger.Info("gather.Start", "sleeping", loopDelay.String())
		time.Sleep(loopDelay)
	}
}

func gather(config *config.App, projects config.Projects) error {
	// Loop through all the projects and grab all of the IPs in the various regions.
	for _, project := range projects.Projects {
		for _, region := range project.Regions {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// TODO: Get CloudSQL IPs.
			// Get NAT IPs.
			ips, err := gnatips.Get(ctx, project.Name, region)
			if err != nil {
				config.Logger.Error(err.Error())
			}
			config.Logger.Info("gather", "project", project.Name, "region", region, "ips", ips)
		}
		// TODO: Store the project in memory.
	}
	return nil
}
