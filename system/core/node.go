package core

import (
	"encoding/json"
	"errors"
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_client"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"sync"
	"time"
)

type Nodes []*Node

type Node struct {
	*m.Node
	mqttClient *mqtt_client.Client
	mqtt       *mqtt.Mqtt
	sync.Mutex
	stat NodeStat
	quit chan struct{}
	ch   map[int64]chan *NodeResponse
}

func NewNode(model *m.Node, mqtt *mqtt.Mqtt) *Node {

	node := &Node{
		Node: model,
		//connStatus: "disabled",
		stat: NodeStat{
			ConnStatus: "disabled",
			LastPing:   time.Now(),
		},
		ch:   make(map[int64]chan *NodeResponse, 0),
		mqtt: mqtt,
		quit: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for {
			time.Sleep(time.Millisecond * 100)
			select {
			case <-ticker.C:
				node.Lock()
				node.stat.IsConnected = time.Now().Sub(node.stat.LastPing).Seconds() < 2

				if node.Node.Status == "enabled" {
					if node.stat.IsConnected {
						node.stat.ConnStatus = "connected"
					} else {
						if time.Now().Sub(node.stat.LastPing).Seconds() < 5 {
							node.stat.ConnStatus = "busy"
						} else {
							node.stat.ConnStatus = "error"
						}
					}
				} else {
					node.stat.ConnStatus = "disabled"
				}

				node.Unlock()
			case <-node.quit:
				close(node.quit)
				return
			default:

			}
		}
	}()

	return node
}

func (n *Node) Remove() {
	n.quit <- struct{}{}
}

func (n *Node) Send(device *m.Device, command []byte) (result NodeResponse, err error) {

	//log.Debugf("send device(%v) command(%v)", device.Id, command)

	// time metric
	startTime := time.Now()

	ch := make(chan *NodeResponse)
	n.addCh(device.Id, ch)
	defer n.delCh(device.Id)

	// send message to node
	msg := &NodeMessage{
		DeviceId:   device.Id,
		DeviceType: device.Type,
		Properties: device.Properties,
		Command:    command,
	}

	n.MqttPublish(msg)

	// wait response
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	var done bool
	for ; ; {
		if done {
			break
		}
		select {
		case <-ticker.C:
			//log.Debugf("request timeout device(%d)", device.Id)
			err = errors.New("timeout")
			done = true
		case resp := <-ch:
			if resp == nil {
				return
			}
			if resp.DeviceId != device.Id {
				continue
			}

			// response from node
			result = NodeResponse{
				DeviceId:   resp.DeviceId,
				Status:     resp.Status,
				DeviceType: resp.DeviceType,
				Response:   resp.Response,
				Properties: resp.Properties,
			}
			done = true
		}
	}

	result.Time = time.Since(startTime).Seconds()

	return
}

func (n *Node) addCh(deviceId int64, ch chan *NodeResponse) {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[deviceId]; ok {
		return
	}

	n.ch[deviceId] = ch
}

func (n *Node) delCh(deviceId int64) {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[deviceId]; !ok {
		return
	}

	close(n.ch[deviceId])
	delete(n.ch, deviceId)
}

func (n *Node) Connect() *Node {

	var err error
	if n.mqttClient, err = n.mqtt.NewClient(nil); err != nil {
		log.Error(err.Error())
	}

	if err = n.mqttClient.Connect(); err != nil {
		log.Error(err.Error())
	}

	time.Sleep(time.Second)

	// /home/node/resp
	if err := n.mqttClient.Subscribe(n.topic("resp"), 0, n.onPublish); err != nil {
		log.Warning(err.Error())
	}

	// /home/node/ping
	if err := n.mqttClient.Subscribe(n.topic("ping"), 0, n.ping); err != nil {
		log.Warning(err.Error())
	}

	return n
}

func (n *Node) Disconnect() {
	if n.mqttClient != nil {
		n.mqttClient.Disconnect()
	}

	_ = n.mqtt.Management().CloseClient(n.Node.Login)
}

func (n *Node) IsConnected() bool {
	n.Lock()
	defer n.Unlock()
	return n.stat.IsConnected
}

func (n *Node) onPublish(client MQTT.Client, msg MQTT.Message) {

	resp := &NodeResponse{}
	if err := json.Unmarshal(msg.Payload(), resp); err != nil {
		log.Error(err.Error())
		return
	}

	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[resp.DeviceId]; !ok {
		return
	}

	n.ch[resp.DeviceId] <- resp
}

func (n *Node) ping(client MQTT.Client, msg MQTT.Message) {

	var stat NodeStatModel
	_ = json.Unmarshal(msg.Payload(), &stat)

	n.Lock()

	//n.stat.Status = stat.Status //????
	n.stat.Thread = stat.Thread
	n.stat.Rps = stat.Rps
	n.stat.Min = stat.Min
	n.stat.Max = stat.Max
	n.stat.StartedAt = stat.StartedAt
	n.stat.LastPing = time.Now()

	n.Unlock()

	return
}

func (n *Node) MqttPublish(msg interface{}) {

	data, _ := json.Marshal(msg)
	if err := n.mqttClient.Publish(n.topic("req"), data); err != nil {
		log.Error(err.Error())
		return
	}
}

func (n *Node) topic(r string) string {
	return fmt.Sprintf("/home/node/%s/%s", n.Node.Name, r)
}

func (n *Node) GetConnStatus() string {
	n.Lock()
	defer n.Unlock()
	return n.stat.ConnStatus
}

func (n *Node) GetStat() NodeStat {
	n.Lock()
	defer n.Unlock()
	return n.stat
}

func (n *Node) UpdateOptions(params *m.Node) {
	n.Lock()
	n.Node = params
	n.Unlock()
}
