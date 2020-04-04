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

package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type AlexaApplications struct {
	Db *gorm.DB
}

type AlexaApplication struct {
	Id            int64 `gorm:"primary_key"`
	ApplicationId string
	Description   string
	Intents       []*AlexaIntent
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (d *AlexaApplication) TableName() string {
	return "alexa_applications"
}

func (n AlexaApplications) Add(v *AlexaApplication) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n AlexaApplications) GetById(id int64) (v *AlexaApplication, err error) {
	v = &AlexaApplication{Id: id}
	err = n.Db.First(&v).
		Preload("Intents").Error
	return
}

func (n *AlexaApplications) List(limit, offset int64, orderBy, sort string) (list []*AlexaApplication, total int64, err error) {

	if err = n.Db.Model(AlexaApplication{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*AlexaApplication, 0)
	q := n.Db.Model(&AlexaApplication{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Preload("Intents").
		Error

	return
}

func (n AlexaApplications) Update(v *AlexaApplication) (err error) {
	err = n.Db.Model(v).Where("id = ?", v.Id).Updates(&map[string]interface{}{
		"application_id": v.ApplicationId,
	}).Error
	return
}

func (n AlexaApplications) Delete(id int64) (err error) {
	err = n.Db.Model(&AlexaApplication{}).Delete("id = ?", id).Error
	return
}
