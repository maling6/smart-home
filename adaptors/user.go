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

package adaptors

import (
	"encoding/json"
	"errors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
	"unicode/utf8"
)

// User ...
type User struct {
	table *db.Users
	db    *gorm.DB
}

// GetUserAdaptor ...
func GetUserAdaptor(d *gorm.DB) *User {
	return &User{
		table: &db.Users{Db: d},
		db:    d,
	}
}

// Add ...
func (n *User) Add(user *m.User) (id int64, err error) {

	dbUser := n.toDb(user)
	dbUser.History.UnmarshalJSON([]byte("[]"))
	if id, err = n.table.Add(dbUser); err != nil {
		return
	}

	metaAdaptor := GetUserMetaAdaptor(n.db)
	for _, meta := range user.Meta {
		meta.UserId = id
		metaAdaptor.UpdateOrCreate(meta)
	}

	return
}

// GetById ...
func (n *User) GetById(userId int64) (user *m.User, err error) {

	var dbUser *db.User
	if dbUser, err = n.table.GetById(userId); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	roleAdaptor := GetRoleAdaptor(n.db)
	err = roleAdaptor.GetAccessList(user.Role)

	return
}

// GetByNickname ...
func (n *User) GetByNickname(nick string) (user *m.User, err error) {

	var dbUser *db.User
	if dbUser, err = n.table.GetByNickname(nick); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	roleAdaptor := GetRoleAdaptor(n.db)
	err = roleAdaptor.GetAccessList(user.Role)

	return
}

// GetByEmail ...
func (n *User) GetByEmail(email string) (user *m.User, err error) {

	var dbUser *db.User
	if dbUser, err = n.table.GetByEmail(email); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	roleAdaptor := GetRoleAdaptor(n.db)
	err = roleAdaptor.GetAccessList(user.Role)

	return
}

// GetByAuthenticationToken ...
func (n *User) GetByAuthenticationToken(token string) (user *m.User, err error) {

	var dbUser *db.User
	if dbUser, err = n.table.GetByAuthenticationToken(token); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	roleAdaptor := GetRoleAdaptor(n.db)
	err = roleAdaptor.GetAccessList(user.Role)

	return
}

// GetByResetPassToken ...
func (n *User) GetByResetPassToken(token string) (user *m.User, err error) {

	if utf8.RuneCountInString(token) > 255 {
		return
	}

	var dbUser *db.User
	if dbUser, err = n.table.GetByResetPassToken(token); err != nil {
		return
	}

	user = n.fromDb(dbUser)

	t := time.Now()
	sub := t.Sub(user.ResetPasswordSentAt.Add(time.Hour * 24)).String()
	if !strings.Contains(sub, "-") {
		err = errors.New("max 24 hour")
	}

	n.ClearResetPassToken(user)

	return
}

// Update ...
func (n *User) Update(user *m.User) (err error) {

	dbUser := n.toDb(user)
	if err = n.table.Update(dbUser); err != nil {
		return
	}

	metaAdaptor := GetUserMetaAdaptor(n.db)
	for _, meta := range user.Meta {
		meta.UserId = user.Id
		metaAdaptor.UpdateOrCreate(meta)
	}

	return
}

// Delete ...
func (n *User) Delete(userId int64) (err error) {
	err = n.table.Delete(userId)
	return
}

// List ...
func (n *User) List(limit, offset int64, orderBy, sort string) (list []*m.User, total int64, err error) {
	var dbList []*db.User
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.User, 0)
	for _, dbUser := range dbList {
		user := n.fromDb(dbUser)
		list = append(list, user)
	}

	return
}

// SignIn ...
func (n *User) SignIn(u *m.User, ipv4 string) (err error) {

	// update count
	u.SignInCount += 1

	// update time
	lastT := u.CurrentSignInAt
	currentT := time.Now()

	u.LastSignInAt = lastT
	u.CurrentSignInAt = &currentT

	// update ipv4
	lastIp := u.CurrentSignInIp
	currentIp := ipv4
	u.LastSignInIp = lastIp
	u.CurrentSignInIp = currentIp

	u.UpdateHistory(currentT, currentIp)

	dbUser := n.toDb(u)
	err = n.table.Update(dbUser)

	return
}

// GenResetPassToken ...
func (n *User) GenResetPassToken(u *m.User) (token string, err error) {

	for {
		token = common.RandStr(50, common.Alphanum)
		u.ResetPasswordToken = token

		if _, err = n.table.GetByResetPassToken(token); err != nil {
			break
		}
	}

	err = n.table.NewResetPassToken(u.Id, u.ResetPasswordToken)

	return
}

// ClearResetPassToken ...
func (n *User) ClearResetPassToken(u *m.User) (err error) {

	err = n.table.ClearResetPassToken(u.Id)
	return
}

// NewToken ...
func (n *User) NewToken(u *m.User) (token string, err error) {

	for {
		token = common.RandStr(50, common.Alphanum)
		u.AuthenticationToken = &token

		if _, err = n.GetByAuthenticationToken(token); err != nil {
			break
		}
	}

	err = n.table.UpdateAuthenticationToken(u.Id, *u.AuthenticationToken)

	return
}

// ClearToken ...
func (n *User) ClearToken(u *m.User) (err error) {

	err = n.table.ClearToken(u.Id)

	return
}

func (n *User) fromDb(dbUser *db.User) (user *m.User) {
	user = &m.User{
		Id:                  dbUser.Id,
		Nickname:            dbUser.Nickname,
		FirstName:           dbUser.FirstName,
		LastName:            dbUser.LastName,
		EncryptedPassword:   dbUser.EncryptedPassword,
		Email:               dbUser.Email,
		Status:              dbUser.Status,
		ResetPasswordToken:  dbUser.ResetPasswordToken,
		AuthenticationToken: dbUser.AuthenticationToken,
		ImageId:             dbUser.ImageId,
		SignInCount:         dbUser.SignInCount,
		CurrentSignInIp:     dbUser.CurrentSignInIp,
		LastSignInIp:        dbUser.LastSignInIp,
		Lang:                dbUser.Lang,
		UserId:              dbUser.UserId,
		RoleName:            dbUser.RoleName,
		ResetPasswordSentAt: dbUser.ResetPasswordSentAt,
		CurrentSignInAt:     dbUser.CurrentSignInAt,
		LastSignInAt:        dbUser.LastSignInAt,
		CreatedAt:           dbUser.CreatedAt,
		UpdatedAt:           dbUser.UpdatedAt,
		DeletedAt:           dbUser.DeletedAt,
		Meta:                make([]*m.UserMeta, 0),
	}

	if dbUser.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		user.Image = imageAdaptor.fromDb(dbUser.Image)
	}

	if dbUser.Meta != nil && len(dbUser.Meta) > 0 {
		userMetaAdaptor := GetUserMetaAdaptor(n.db)
		for _, dbMeta := range dbUser.Meta {
			meta := userMetaAdaptor.fromDb(dbMeta)
			user.Meta = append(user.Meta, meta)
		}
	}

	// history
	user.History = make([]*m.UserHistory, 0)
	data, _ := dbUser.History.MarshalJSON()
	json.Unmarshal(data, &user.History)

	// role
	if dbUser.Role != nil {
		roleAdaptor := GetRoleAdaptor(n.db)
		user.Role = roleAdaptor.fromDb(dbUser.Role)
	}

	// created by
	if dbUser.User != nil {
		user.User = n.fromDb(dbUser.User)
	}

	return
}

func (n *User) toDb(user *m.User) (dbUser *db.User) {
	dbUser = &db.User{
		Id:                  user.Id,
		Nickname:            user.Nickname,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		EncryptedPassword:   user.EncryptedPassword,
		Email:               user.Email,
		Status:              user.Status,
		ResetPasswordToken:  user.ResetPasswordToken,
		AuthenticationToken: user.AuthenticationToken,
		ImageId:             user.ImageId,
		SignInCount:         user.SignInCount,
		CurrentSignInIp:     user.CurrentSignInIp,
		LastSignInIp:        user.LastSignInIp,
		Lang:                user.Lang,
		UserId:              user.UserId,
		RoleName:            user.RoleName,
		ResetPasswordSentAt: user.ResetPasswordSentAt,
		CurrentSignInAt:     user.CurrentSignInAt,
		LastSignInAt:        user.LastSignInAt,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
		DeletedAt:           user.DeletedAt,
	}

	if user.ImageId.Valid {
		dbUser.ImageId.Scan(user.ImageId.Int64)
	}

	if user.UserId.Valid {
		dbUser.UserId.Scan(user.UserId.Int64)
	}

	return
}
