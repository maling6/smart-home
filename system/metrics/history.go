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

package metrics

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/gammazero/deque"
	"sync"
	"time"
)

type History struct {
	Items []HistoryItem `json:"items"`
}

type HistoryManager struct {
	publisher  IPublisher
	adaptors   *adaptors.Adaptors
	updateLock sync.Mutex
	queue      *deque.Deque
}

func NewHistoryManager(publisher IPublisher, adaptors *adaptors.Adaptors) *HistoryManager {
	manager := &HistoryManager{
		publisher: publisher,
		adaptors:  adaptors,
		queue:     &deque.Deque{},
	}
	manager.init()
	return manager
}

func (d *HistoryManager) update(t interface{}) {
	switch v := t.(type) {
	case HistoryItem:
		d.updatePool(v)
	default:
		return
	}

	d.broadcast()
}

func (d *HistoryManager) updatePool(item HistoryItem) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	const (
		maxItems = 8
	)

	for d.queue.Len() > maxItems {
		d.queue.PopBack()
	}

	d.queue.PushFront(HistoryItem{
		DeviceName:        item.DeviceName,
		DeviceDescription: item.DeviceDescription,
		Type:              item.Type,
		Description:       item.Description,
		CreatedAt:         item.CreatedAt,
	})
}

func (d HistoryManager) Snapshot() History {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	items := make([]HistoryItem, d.queue.Len())
	for i := 0; i < d.queue.Len(); i++ {
		if item, ok := d.queue.At(i).(HistoryItem); ok {
			items[i] = item
		}
	}

	return History{
		Items: items,
	}
}

func (d *HistoryManager) broadcast() {
	go d.publisher.Broadcast("history")
}

func (d *HistoryManager) init() {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	if list, err := d.adaptors.MapDeviceHistory.List(8, 0); err == nil {
		for _, item := range list {
			d.queue.PushBack(HistoryItem{
				DeviceName:        item.MapElement.Name,
				DeviceDescription: item.MapElement.Description,
				Type:              string(item.Type),
				Description:       item.Description,
				CreatedAt:         item.CreatedAt,
			})
		}
	}

}

type HistoryItem struct {
	DeviceName        string    `json:"device_name"`
	DeviceDescription string    `json:"device_description"`
	Type              string    `json:"type"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"created_at"`
}
