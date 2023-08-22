package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/darron/gips/config"
	"github.com/darron/gips/service"
	"github.com/spf13/cobra"
)

var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "HTTP Service Commands",
		Run: func(cmd *cobra.Command, args []string) {
			StartService()
		},
	}

	defaultStorageLayer = "memory"
	storageLayer        string

	defaultMiddlewareHTMLCacheEnabled = true
	middlewareHTMLCacheEnabled        bool
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVarP(&storageLayer, "storage", "", GetENVVariable("STORAGE", defaultStorageLayer), "Storage Layer: memory")
	serviceCmd.Flags().BoolVarP(&middlewareHTMLCacheEnabled, "htmlcache", "", GetBoolENVVariable("HTMLCACHE_ENABLED", defaultMiddlewareHTMLCacheEnabled), "Enable Middleware Cache")
}

func StartService() {
	// Setup some options
	var opts []config.OptFunc
	var tlsConfig *tls.Config

	opts = append(opts, config.WithPort(port))
	opts = append(opts, config.WithLogger(logLevel, logFormat))
	opts = append(opts, config.WithMiddlewareHTMLCache(middlewareHTMLCacheEnabled))

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
