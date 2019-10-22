package password

import (
	"time"

	"github.com/labstack/echo"

	password "github.com/Wallruzz9114/plethora_api/pkg/api/password"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
)

const name = "password"

// LogService represents password logging service
type LogService struct {
	password.Service
	logger plethora_api.Logger
}

// New creates new password logging service
func New(svc password.Service, logger plethora_api.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// Change logging
func (ls *LogService) Change(c echo.Context, id int, oldPass, newPass string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Change password request", err,
			map[string]interface{}{
				"req":  id,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.Service.Change(c, id, oldPass, newPass)
}
