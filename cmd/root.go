package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	binaryName = "gips"
	rootCmd    = &cobra.Command{
		Use:   binaryName,
		Short: binaryName + " collects GCP NAT IPs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().Help() //nolint
		},
	}
	defaultLogLevel = "debug"
	logLevel        string

	defaultLogFormat = "text"
	logFormat        string

	defaultPort = "8080"
	port        string

	defaultEnableTLS            = false
	defaultEnableTLSLetsEncrypt = false
	enableTLS                   bool
	enableTLSLetsEncrypt        bool
	tlsCache                    string
	tlsCert                     string
	tlsDomains                  string
	tlsEmail                    string
	tlsKey                      string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "level", "", GetENVVariable("LOG_LEVEL", defaultLogLevel), "Log level: info, debug, etc")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "format", "", GetENVVariable("LOG_FORMAT", defaultLogFormat), "Log format: text or json")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", GetENVVariable("PORT", defaultPort), "HTTP Port to listen on")
	rootCmd.PersistentFlags().BoolVarP(&enableTLS, "tls", "", GetBoolENVVariable("ENABLE_TLS", defaultEnableTLS), "Enable TLS")
	rootCmd.PersistentFlags().BoolVarP(&enableTLSLetsEncrypt, "letsencrypt", "", GetBoolENVVariable("ENABLE_TLS_LETS_ENCRYPT", defaultEnableTLSLetsEncrypt), "Enable LetsEncrypt")
	rootCmd.PersistentFlags().StringVarP(&tlsCache, "tlsCache", "", GetENVVariable("TLS_CACHE", ""), "Cache Dir for TLS Certificate")
	rootCmd.PersistentFlags().StringVarP(&tlsCert, "cert", "", GetENVVariable("TLS_CERT", ""), "Manual TLS Cert")
	rootCmd.PersistentFlags().StringVarP(&tlsDomains, "tlsDomains", "", GetENVVariable("TLS_DOMAINS", ""), "Domains for TLS Certificate - separate by commas")
	rootCmd.PersistentFlags().StringVarP(&tlsEmail, "tlsEmail", "", GetENVVariable("TLS_EMAIL", ""), "Email Address for TLS Certificate")
	rootCmd.PersistentFlags().StringVarP(&tlsKey, "key", "", GetENVVariable("TLS_KEY", ""), "Manual TLS Key")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetENVVariable(name, defValue string) string {
	v, ok := os.LookupEnv(name)
	if ok {
		return v
	}
	return defValue
}

func GetBoolENVVariable(name string, defValue bool) bool {
	_, ok := os.LookupEnv(name)
	if ok {
		return true
	}
	return defValue
}
