package auth

import (
	"time"

	"github.com/labstack/echo"

	auth "github.com/Wallruzz9114/plethora_api/pkg/api/auth"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
)

const name = "auth"

// LogService represents auth logging service
type LogService struct {
	auth.Service
	logger plethora_api.Logger
}

// New creates new auth logging service
func New(svc auth.Service, logger plethora_api.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// Authenticate logging
func (ls *LogService) Authenticate(c echo.Context, user, password string) (resp *plethora_api.AuthToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Authenticate request", err,
			map[string]interface{}{
				"req":  user,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.Service.Authenticate(c, user, password)
}

// Refresh logging
func (ls *LogService) Refresh(c echo.Context, req string) (resp *plethora_api.RefreshToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Refresh request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.Service.Refresh(c, req)
}

// Me logging
func (ls *LogService) Me(c echo.Context) (resp *plethora_api.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Me request", err,
			map[string]interface{}{
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.Service.Me(c)
}
