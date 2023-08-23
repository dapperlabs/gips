package gather

import (
	"fmt"
	"time"

	"github.com/darron/gips/config"
)

var (
	loopDelay = time.Minute * 5
)

func Start(config *config.App, projects config.Projects) {
	// Start gathering data from GCP.
	// Loop forever.
	for {
		config.Logger.Info("Gathering data from GCP")
		err := gather(config, projects)
		if err != nil {
			config.Logger.Error(err.Error())
		}
		config.Logger.Info(fmt.Sprintf("Sleeping for %s", loopDelay))
		time.Sleep(loopDelay)
	}
}

func gather(config *config.App, projects config.Projects) error {
	config.Logger.Info("into gather")
	return nil
}
