package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/darron/gips/config"
	"github.com/darron/gips/gather"
	"github.com/darron/gips/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "HTTP Service Commands",
		Run: func(cmd *cobra.Command, args []string) {
			StartService()
		},
	}

	defaultConfigFileName = "config"
	configFileName        string

	defaultConfigFileExtension = "yaml"
	configFileExtension        string

	defaultStorageLayer = "memory"
	storageLayer        string
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVarP(&configFileName, "config", "", GetENVVariable("CONFIG_FILE_NAME", defaultConfigFileName), "Config file name (without extension)")
	serviceCmd.Flags().StringVarP(&configFileExtension, "fileExtension", "", GetENVVariable("CONFIG_FILE_EXTENSION", defaultConfigFileExtension), "Config file extension")
	serviceCmd.Flags().StringVarP(&storageLayer, "storage", "", GetENVVariable("STORAGE", defaultStorageLayer), "Storage Layer: memory")
}

func StartService() {
	// We need a config file for all the projects and regions.
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}
	// Let's put it into the config.Projects struct.
	projects := config.Projects{}
	err = viper.Unmarshal(&projects)
	if err != nil {
		log.Fatal("Config Unmarshal Error: ", err)
	}

	for _, project := range projects.Projects {
		log.Printf("Project: %q in %s\n", project.Name, project.Regions)
	}

	// Setup some options
	var opts []config.OptFunc
	var tlsConfig *tls.Config

	// Let's setup the storage layer
	switch storageLayer {
	case "memory":
		opts = append(opts, config.WithMemoryStore())
	default:
		log.Fatal("Unknown storage layer")
	}

	opts = append(opts, config.WithPort(port))
	opts = append(opts, config.WithLogger(logLevel, logFormat))

	// Let's turn on TLS with LetsEncrypt
	// Setup the config here.
	if enableTLS && enableTLSLetsEncrypt {
		log.Println("Enabling LetsEncrypt")
		tlsVar := config.TLS{
			CacheDir:    tlsCache,
			DomainNames: tlsDomains,
			Email:       tlsEmail,
			Enable:      enableTLS,
		}
		err := tlsVar.LetsEncryptVerify()
		if err != nil {
			log.Fatal(err)
		}
		// Let's setup the service http.Server tls.Config
		tlsConfig = tlsVar.LetsEncryptTLSConfig()
		opts = append(opts, config.WithTLS(tlsVar))
	}

	// If we have manually generated certs - let's use those for HTTPS
	// Setup the config here.
	if enableTLS && !enableTLSLetsEncrypt && (tlsCert != "") && (tlsKey != "") {
		log.Println("Enabling TLS with manual certs")
		tlsVar := config.TLS{
			CertFile:    tlsCert,
			DomainNames: tlsDomains,
			Enable:      enableTLS,
			KeyFile:     tlsKey,
		}
		err := tlsVar.StaticCredentialsVerify()
		if err != nil {
			log.Fatal(err)
		}
		tlsConfig, err = tlsVar.StaticCredentialsTLSConfig()
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, config.WithTLS(tlsVar))
	}

	// Let's get the config for the app
	conf, err := config.Get(opts...)
	if err != nil {
		log.Fatal(err)
	}

	// Start gathering data from GCP.
	go gather.Start(conf, projects)

	conf.Logger.Info("Starting HTTP Service")
	s, err := service.Get(conf)
	if err != nil {
		conf.Logger.Error(err.Error())
		os.Exit(1)
	}

	// If we are going to turn on TLS - let's launch it.
	if enableTLS {
		h := http.Server{
			Addr:        ":443",
			Handler:     s,
			TLSConfig:   tlsConfig,
			ReadTimeout: 30 * time.Second, // use custom timeouts
		}
		if err := h.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			s.Logger.Fatal(err)
		}
	}
	s.Logger.Fatal(s.Start(":" + conf.Port))
}
