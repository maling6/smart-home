package models

type MapOptions struct {
	Zoom              float64 `json:"zoom"`
	ElementStateText  bool    `json:"element_state_text"`
	ElementOptionText bool    `json:"element_option_text"`
}

// swagger:model
type NewMap struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
}

// swagger:model
type UpdateMap struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
}

// swagger:model
type Map struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
}

// swagger:model
type MapFull struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Options     MapOptions  `json:"options"`
	Layers      []*MapLayer `json:"layers"`
}