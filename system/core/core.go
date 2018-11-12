package core

import (
	"sync"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/op/go-logging"
	"github.com/e154/smart-home/system/scripts"
)

var (
	log = logging.MustGetLogger("core")
)

type Core struct {
	sync.Mutex
	nodes     map[int64]*m.Node
	workflows map[int64]*Workflow
	adaptors  *adaptors.Adaptors
	scripts   *scripts.ScriptService
	//telemetry		Telemetry
	//Map				*Map
}

func NewCore(adaptors *adaptors.Adaptors,
	scripts *scripts.ScriptService) (core *Core, err error) {
	core = &Core{
		nodes:     make(map[int64]*m.Node),
		workflows: make(map[int64]*Workflow),
		adaptors:  adaptors,
		scripts:   scripts,
	}

	return
}

func (c *Core) Run() (err error) {
	if err = c.initNodes(); err != nil {
		return
	}
	err = c.InitWorkflows()
	return
}

func (b *Core) Stop() (err error) {

	for _, workflow := range b.workflows {
		if err = b.DeleteWorkflow(workflow.model); err != nil {
			return
		}
	}

	for _, node := range b.nodes {
		if err = b.RemoveNode(node); err != nil {
			return
		}
	}

	return
}

func (b *Core) Restart() (err error) {

	if err = b.Stop(); err != nil {
		log.Error(err.Error())
	}

	if err = b.Run(); err != nil {
		log.Error(err.Error())
	}

	return
}

// ------------------------------------------------
// Nodes
// ------------------------------------------------

func (c *Core) initNodes() (err error) {

	var nodes []*m.Node
	if nodes, err = c.adaptors.Node.GetAllEnabled(); err != nil {
		return
	}

	for _, model_node := range nodes {
		c.AddNode(model_node)
	}

	return
}

func (b *Core) AddNode(node *m.Node) (err error) {

	if _, exist := b.nodes[node.Id]; exist {
		return b.ReloadNode(node)
	}

	log.Infof("Add node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		return
	}

	b.Lock()
	node.Start()
	b.nodes[node.Id] = node.Connect()
	b.Unlock()

	//TODO add telemetry
	//b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) RemoveNode(node *m.Node) (err error) {

	log.Infof("Remove node: \"%s\"", node.Name)

	if _, exist := b.nodes[node.Id]; !exist {
		return
	}

	b.Lock()
	if _, ok := b.nodes[node.Id]; ok {
		b.nodes[node.Id].Disconnect()
		delete(b.nodes, node.Id)
	}

	delete(b.nodes, node.Id)
	b.Unlock()

	//TODO add telemetry
	//b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) ReloadNode(node *m.Node) (err error) {

	log.Infof("Reload node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; !ok {
		b.AddNode(node)
		return
	}

	b.Lock()
	b.nodes[node.Id].Status = node.Status
	b.nodes[node.Id].Ip = node.Ip
	b.nodes[node.Id].Port = node.Port
	b.nodes[node.Id].SetConnectStatus("wait")
	b.Unlock()

	if b.nodes[node.Id].Status == "disabled" {
		b.nodes[node.Id].Disconnect()
	} else {
		b.nodes[node.Id].Connect()
	}

	return
}

func (b *Core) ConnectNode(node *m.Node) (err error) {

	log.Infof("Connect to node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes[node.Id].Connect()
	}

	//TODO add telemetry
	//b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) DisconnectNode(node *m.Node) (err error) {

	log.Infof("Disconnect from node: \"%s\"", node.Name)

	if _, ok := b.nodes[node.Id]; ok {
		b.nodes[node.Id].Disconnect()
	}

	//TODO add telemetry
	//b.telemetry.Broadcast("nodes")

	return
}

func (b *Core) GetNodes() (nodes map[int64]*m.Node) {

	nodes = make(map[int64]*m.Node)

	b.Lock()
	for id, node := range b.nodes {
		nodes[id] = node
	}
	b.Unlock()

	return
}

// ------------------------------------------------
// Workflows
// ------------------------------------------------

// инициализация всего рабочего процесса, с запуском
// дочерни подпроцессов
func (b *Core) InitWorkflows() (err error) {

	workflows, err := b.adaptors.Workflow.GetAllEnabled()
	if err != nil {
		return
	}

	for _, workflow := range workflows {
		if err = b.AddWorkflow(workflow); err != nil {
			return
		}
	}

	return
}

// добавление рабочего процесс
func (b *Core) AddWorkflow(workflow *m.Workflow) (err error) {

	log.Infof("Add workflow: %s", workflow.Name)

	if _, ok := b.workflows[workflow.Id]; ok {
		return
	}

	wf := NewWorkflow(workflow, b.adaptors, b.scripts)

	if err = wf.Run(); err != nil {
		return
	}

	b.workflows[workflow.Id] = wf

	return
}

// нельзя удалить workflow, если присутствуют связанные сущности
func (b *Core) DeleteWorkflow(workflow *m.Workflow) (err error) {

	log.Infof("Remove workflow: %s", workflow.Name)

	if _, ok := b.workflows[workflow.Id]; !ok {
		return
	}

	b.workflows[workflow.Id].Stop()
	delete(b.workflows, workflow.Id)

	return
}

func (b *Core) UpdateWorkflowScenario(workflow *m.Workflow) (err error) {

	if _, ok := b.workflows[workflow.Id]; !ok {
		return
	}

	err = b.workflows[workflow.Id].UpdateScenario()

	return
}

func (b *Core) UpdateWorkflow(workflow *m.Workflow) (err error) {

	if workflow.Status == "enabled" {
		if _, ok := b.workflows[workflow.Id]; !ok {
			err = b.AddWorkflow(workflow)
		}
	} else {
		if _, ok := b.workflows[workflow.Id]; ok {
			err = b.DeleteWorkflow(workflow)
		}
	}

	return
}