package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/darron/gips/adaptors/memory"
	"github.com/darron/gips/core"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type OptFunc func(*Opts)

type Opts struct {
	Logger              *slog.Logger
	MiddlewareHTMLCache bool
	Port                string
	ProjectStore        core.ProjectService
	TLS                 TLS
}

type TLS struct {
	CacheDir    string
	CertFile    string
	DomainNames string
	Email       string
	Enable      bool
	KeyFile     string
}

type App struct {
	Opts
}

type Projects struct {
	Projects []Project
}

type Project struct {
	Name    string
	Regions []string
}

var (
	defaultLogformat           = "text"
	defaultLogLevel            = "debug"
	defaultPort                = "8080"
	defaultHTMLMiddlewareCache = true
)

func defaultOpts() Opts {
	return Opts{
		Port:                defaultPort,
		MiddlewareHTMLCache: true,
	}
}

func WithMemoryStore() OptFunc {
	return func(opts *Opts) {
		opts.ProjectStore = memory.New()
	}
}

func WithMiddlewareHTMLCache(enabled bool) OptFunc {
	return func(opts *Opts) {
		opts.MiddlewareHTMLCache = enabled
	}
}

func WithLogger(level, format string) OptFunc {
	l := GetLogger(level, format)
	return func(opts *Opts) {
		opts.Logger = l
	}
}

func WithPort(port string) OptFunc {
	return func(opts *Opts) {
		opts.Port = port
	}
}

func WithTLS(tls TLS) OptFunc {
	return func(opts *Opts) {
		opts.TLS = tls
	}
}

func New() (*App, error) {
	var optFuncs []OptFunc

	// Really basic default options without any configuration.
	// Moving configuration to cmd and will be calling `Get(WithOption())`
	optFuncs = append(optFuncs, WithLogger(defaultLogLevel, defaultLogformat))
	optFuncs = append(optFuncs, WithPort(defaultPort))
	optFuncs = append(optFuncs, WithMiddlewareHTMLCache(defaultHTMLMiddlewareCache))

	return Get(optFuncs...)
}

func Get(opts ...OptFunc) (*App, error) {
	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}
	app := App{
		Opts: o,
	}
	return &app, nil
}

func (a *App) GetHTTPEndpoint() string {
	protocol := "http"
	domain := "localhost"
	port := a.Port
	if a.TLS.DomainNames != "" {
		protocol = "https"
		domain = strings.Split(a.TLS.DomainNames, ",")[0]
		port = "443"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, domain, port)
}

func GetLogger(level, format string) *slog.Logger {
	var slogLevel slog.Level
	var slogHandler slog.Handler

	// Let's deal with level.
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	default:
		slogLevel = slog.LevelInfo
	}
	handlerOpts := slog.HandlerOptions{
		Level: slogLevel,
	}

	// Let's switch formats as desired.
	switch format {
	case "json":
		slogHandler = slog.NewJSONHandler(os.Stdout, &handlerOpts)
	default:
		slogHandler = slog.NewTextHandler(os.Stdout, &handlerOpts)
	}
	log := slog.New(slogHandler)

	return log
}

func (t TLS) LetsEncryptVerify() error {
	if t.CacheDir == "" {
		return errors.New("Cache dir cannot be emtpy")
	}
	// Check to see if the cache dir exists - if it doesn't try to create it.
	if _, err := os.Open(t.CacheDir); os.IsNotExist(err) {
		// It doesn't exist - try to create it.
		err := os.MkdirAll(t.CacheDir, 0750)
		if err != nil {
			return err
		}
	}
	if t.DomainNames == "" {
		return errors.New("Domain names cannot be empty")
	}
	if t.Email == "" {
		return errors.New("Email address cannot be empty")
	}
	return nil
}

func (t TLS) StaticCredentialsVerify() error {
	_, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return err
	}
	return nil
}

func (t TLS) StaticCredentialsTLSConfig() (*tls.Config, error) {
	var tlsConfig *tls.Config
	cer, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return tlsConfig, err
	}
	tlsConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
	return tlsConfig, nil
}

func (t TLS) LetsEncryptTLSConfig() *tls.Config {
	domains := strings.Split(t.DomainNames, ",")
	autoTLSManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(t.CacheDir),
		Email:      t.Email,
		HostPolicy: autocert.HostWhitelist(domains...),
	}
	tlsConfig := tls.Config{
		GetCertificate: autoTLSManager.GetCertificate,
		NextProtos:     []string{acme.ALPNProto},
	}
	return &tlsConfig
}
