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

package core

// Javascript Binding
//
// Action
//	 .Id
//	 .Name
//	 .Description
//	 .Device()
//
type ActionBind struct {
	Id          int64
	Name        string
	Description string
	action      *Action
}

// NewActionBind ...
func NewActionBind(id int64, name, desc string, action *Action) *ActionBind {
	return &ActionBind{
		Id:          id,
		Name:        name,
		Description: desc,
		action:      action,
	}
}

// Device ...
func (a *ActionBind) Device() *DeviceBind {
	return &DeviceBind{model: a.action.GetDevice()}
}

//func (a *ActionBind) Node() *NodeBind {
//	return &NodeBind{node: a.action.GetNode()}
//}
