package service

import (
	"errors"
	"net/http"
	"net/netip"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	ErrorInvalidIP    = errors.New("not a valid IP address")
	ErrorWellFormedIP = errors.New("not a well-formed or valid IP address")
	ErrorNotIPv4      = errors.New("not a valid IPv4 address")
	ErrorPrivateIP    = errors.New("not a routable IP address")
	ErrorLoopbackIP   = errors.New("that's a loopback IP address")
)

func (s HTTPService) IPSearch(c echo.Context) error {
	ip := c.Param("ip")
	if ip == "" {
		return c.JSON(http.StatusNotFound, "ip must not be blank")
	}
	// Check to see if it's well formed and valid.
	err := isIPValid(ip)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Now let's see if we can find it.
	outputProject, err := s.conf.ProjectStore.FindIP(ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if outputProject == nil {
		return c.JSON(http.StatusNotFound, "no project found")
	}
	return c.JSON(http.StatusOK, *outputProject)
}

func isIPValid(ip string) error {
	// Can't be 0.0.0.0
	if strings.HasPrefix(ip, "0.0.0.0") {
		return ErrorInvalidIP
	}
	// Is it well formed and valid?
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ErrorWellFormedIP
	}
	// Only IPv4 for now.
	if !addr.Is4() {
		return ErrorNotIPv4
	}
	// Must be public and routable - can't be a private IP address.
	if addr.IsPrivate() {
		return ErrorPrivateIP
	}
	// Can't be loopback:
	if addr.IsLoopback() {
		return ErrorLoopbackIP
	}
	return nil
}
