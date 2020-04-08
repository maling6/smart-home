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

import "time"

// DeviceDefault ...
type DeviceDefault struct {
	Address *int `json:"address"`
}

// DeviceZigBee ...
type DeviceZigBee struct {
	Address string `json:"address"`
}

// DeviceSmartBus ...
type DeviceSmartBus struct {
	Address  *int          `json:"address"`
	Baud     int           `json:"baud"`
	Sleep    int64         `json:"sleep"`
	StopBite int64         `json:"stop_bite"`
	Timeout  time.Duration `json:"timeout"`
	Tty      string        `json:"tty"`
}
