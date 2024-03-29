package transport

import (
	"net/http"

	"github.com/labstack/echo"

	password "github.com/Wallruzz9114/plethora_api/pkg/api/password"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

// Custom errors
var (
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
)

// Password change request
// swagger:model pwChange
type changeReq struct {
	ID                 int    `json:"-"`
	OldPassword        string `json:"old_password" validate:"required,min=8"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required"`
}

// HTTP represents password http transport service
type HTTP struct {
	svc password.Service
}

// NewHTTP creates new password http service
func NewHTTP(svc password.Service, er *echo.Group) {
	h := HTTP{svc}
	pr := er.Group("/password")

	// swagger:operation PATCH /v1/password/{id} password pwChange
	// ---
	// summary: Changes user's password.
	// description: If user's old passowrd is correct, it will be replaced with new password.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: int
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/pwChange"
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	pr.PATCH("/:id", h.change)
}

func (h *HTTP) change(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))

	if err != nil {
		return plethora_api.ErrBadRequest
	}

	p := new(changeReq)

	if err := c.Bind(p); err != nil {
		return err
	}

	if p.NewPassword != p.NewPasswordConfirm {
		return ErrPasswordsNotMaching
	}

	if err := h.svc.Change(c, id, p.OldPassword, p.NewPassword); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
