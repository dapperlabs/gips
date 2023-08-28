package cmd

import (
	"fmt"
	"log"

	"github.com/darron/gips/config"
	"github.com/darron/gips/search"
	"github.com/spf13/cobra"
)

var (
	searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search " + binaryName + " for an IP address",
		Run: func(cmd *cobra.Command, args []string) {
			doSearch()
		},
	}
	ip string
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&ip, "ip", "", "", "IP Address to Search")
}

func doSearch() {
	if ip == "" {
		log.Fatal("You must provide an IP address to search for.")
	}
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	serviceAddr := conf.GetHTTPEndpoint()
	endpoint := fmt.Sprintf("%s/api/v1/search/%s", serviceAddr, ip)
	fmt.Printf("Searching for %s in %s\n", ip, endpoint)
	project, err := search.IP(ip, endpoint)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("IP address %s not found.\n", ip)
	}
	if project.Name != "" {
		fmt.Printf("Project: %#v\n", project)
	}
}
