package middleware

import (
	"net"
	"net/http"

	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/gofiber/fiber/v2"
)

// мидлвар, который проверяет находится ли ip-адрес клиента
// в доверенной подсети из конфига TrustedSubnet
func IPInTrustedSubnet(ctx *fiber.Ctx) error {
	ip := ctx.Get("X-Real-IP")
	if ip == "" {
		return ctx.SendStatus(http.StatusForbidden)
	}

	cfg := config.New()

	_, subnet, err := net.ParseCIDR(cfg.TrustedSubnet)
	if err != nil {
		logger.Log.Error("could not parse cidr in IPInTrustedSubnet", err.Error())
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if subnet == nil || !subnet.Contains(net.ParseIP(ip)) {
		return ctx.SendStatus(http.StatusForbidden)
	}

	return ctx.Next()
}
