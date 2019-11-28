package models

import "time"

// swagger:model
type MqttClient struct {
	ClientID       string    `json:"client_id"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	KeepAlive      uint16    `json:"keep_alive"`
	CleanSession   bool      `json:"clean_session"`
	WillFlag       bool      `json:"will_flag"`
	WillRetain     bool      `json:"will_retain"`
	WillQos        uint8     `json:"will_qos"`
	WillTopic      string    `json:"will_topic"`
	WillPayload    string    `json:"will_payload"`
	RemoteAddr     string    `json:"remote_addr"`
	LocalAddr      string    `json:"local_addr"`
	ConnectedAt    time.Time `json:"connected_at"`
	DisconnectedAt time.Time `json:"disconnected_at"`
}

// swagger:model
type MqttSession struct {
	ClientID          string    `json:"client_id"`
	Status            string    `json:"status"`
	CleanSession      bool      `json:"clean_session"`
	Subscriptions     int64     `json:"subscriptions"`
	MaxInflight       int       `json:"max_inflight"`
	InflightLen       int64     `json:"inflight_len"`
	MaxMsgQueue       int       `json:"max_msg_queue"`
	MsgQueueLen       int64     `json:"msg_queue_len"`
	MaxAwaitRel       int       `json:"max_await_rel"`
	AwaitRelLen       int64     `json:"await_rel_len"`
	MsgDroppedTotal   int64     `json:"msg_dropped_total"`
	MsgDeliveredTotal int64     `json:"msg_delivered_total"`
	ConnectedAt       time.Time `json:"connected_at"`
	DisconnectedAt    time.Time `json:"disconnected_at"`
}

// swagger:model
type MqttSubscription struct {
	ClientID string    `json:"client_id"`
	Qos      uint8     `json:"qos"`
	Name     string    `json:"name"`
	At       time.Time `json:"at"`
}

// swagger:model
type NewMqttPublish struct {
	Topic   string `json:"topic"`
	Qos     int    `json:"qos"`
	Payload []byte `json:"payload"`
	Retain  bool   `json:"retain"`
}
