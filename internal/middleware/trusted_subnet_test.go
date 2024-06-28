package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPInTrustedSubnet(t *testing.T) {
	app := fiber.New()

	cfg := config.New()

	cfg.TrustedSubnet = "192.168.0.0/24"

	app.Use(IPInTrustedSubnet)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "192.168.0.1")

	resp, err := app.Test(req, 1)

	require.NoError(t, err)

	assert.NotEqual(t, http.StatusForbidden, resp.StatusCode)
}

func TestIPInTrustedSubnetWithIpNotInSubnet(t *testing.T) {
	app := fiber.New()

	cfg := config.New()

	cfg.TrustedSubnet = "145.132.0.0/24"

	app.Use(IPInTrustedSubnet)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "192.168.0.1")

	resp, err := app.Test(req, 1)

	require.NoError(t, err)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestIPInTrustedSubnetWithEmptyIp(t *testing.T) {
	app := fiber.New()

	app.Use(IPInTrustedSubnet)

	req := httptest.NewRequest("GET", "/", nil)

	resp, err := app.Test(req, 1)

	require.NoError(t, err)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
