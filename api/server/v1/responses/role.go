package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response RoleList
type RoleList struct {
	// in:body
	Body struct {
		Items []*models.Role `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response RoleSearch
type RoleSearch struct {
	// in:body
	Body struct {
		Roles []*models.Role `json:"roles"`
	}
}
