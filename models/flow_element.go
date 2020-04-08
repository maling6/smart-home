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

package models

import (
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/uuid"
	"github.com/e154/smart-home/system/validation"
	"time"
)

// FlowElementGraphSettingsPosition ...
type FlowElementGraphSettingsPosition struct {
	Top  int64 `json:"top"`
	Left int64 `json:"left"`
}

// FlowElementGraphSettings ...
type FlowElementGraphSettings struct {
	Position FlowElementGraphSettingsPosition `json:"position"`
}

// FlowElement ...
type FlowElement struct {
	Uuid          uuid.UUID                 `json:"uuid"`
	Name          string                    `json:"name" valid:"MaxSize(254)"`
	Description   string                    `json:"description"`
	FlowId        int64                     `json:"flow_id" valid:"Required"`
	Script        *Script                   `json:"script"`
	ScriptId      *int64                    `json:"script_id"`
	Status        StatusType                `json:"status" valid:"Required"`
	FlowLink      *int64                    `json:"flow_link"`
	PrototypeType FlowElementsPrototypeType `json:"prototype_type" valid:"Required"`
	GraphSettings FlowElementGraphSettings  `json:"graph_settings"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

// Valid ...
func (d *FlowElement) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
