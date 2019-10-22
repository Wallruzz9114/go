package transport

import (
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
)

// Login request
// swagger:parameters login
type swaggLoginReq struct {
	// in:body
	Body credentials
}

// Login response
// swagger:response loginResp
type swaggLoginResp struct {
	// in:body
	Body struct {
		*plethora_api.AuthToken
	}
}

// Token refresh response
// swagger:response refreshResp
type swaggRefreshResp struct {
	// in:body
	Body struct {
		*plethora_api.RefreshToken
	}
}
