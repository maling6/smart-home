// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response TemplateList
type TemplateList struct {
	// in:body
	Body struct {
		Items []*models.Template `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"objects_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response TemplateSearch
type TemplateSearch struct {
	// in:body
	Body struct {
		Templates []*models.Template `json:"templates"`
	}
}
