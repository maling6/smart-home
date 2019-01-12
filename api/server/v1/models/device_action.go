package models

import "time"

type DeviceActionScript struct {
	Id int64 `json:"id"`
}
type DeviceActionDevice struct {
	Id int64 `json:"id"`
}
type NewDeviceAction struct {
	Name        string             `json:"name" valid:"MaxSize(254);Required"`
	Description string             `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	Script      *DeviceActionScript `json:"script"`
}

type UpdateDeviceAction struct {
	Name        string             `json:"name" valid:"MaxSize(254);Required"`
	Description string             `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	Script      *DeviceActionScript `json:"script"`
}

type DeviceAction struct {
	Id          int64              `json:"id"`
	Name        string             `json:"name" valid:"MaxSize(254);Required"`
	Description string             `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	Script      *DeviceActionScript `json:"script"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type DeviceActionListModel struct {
	Items []DeviceAction
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}
