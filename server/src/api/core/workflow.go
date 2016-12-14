package core

import (
	"../log"
	"sync"
	"../models"
)

func NewWorkflow(model *models.Workflow, nodes map[int64]*models.Node) (workflow *Workflow) {

	workflow = &Workflow{
		model: model,
		Nodes: nodes,
		Flows: make(map[int64]*Flow),
		mutex: &sync.Mutex{},
	}

	return
}

type Workflow struct {
	model   	*models.Workflow
	Nodes   	map[int64]*models.Node
	mutex   	*sync.Mutex
	Flows   	map[int64]*Flow
}

func (wf *Workflow) Run() (err error) {

	err = wf.InitFlows()

	if err != nil {
		return
	}

	return
}

func (wf *Workflow) Stop() (err error) {

	for _, flow := range wf.Flows {
		wf.RemoveFlow(flow.Model)
	}

	return
}

func (wf *Workflow) Restart() (err error) {

	wf.Stop()
	err = wf.Run()

	return
}

// ------------------------------------------------
// Flows
// ------------------------------------------------

// получаем все связанные процессы
func (wf *Workflow) InitFlows() (err error) {

	var flows []*models.Flow
	if flows, err = wf.model.GetAllEnabledFlows(); err != nil {
		return
	}

	for _, flow := range flows {
		wf.AddFlow(flow)
	}

	return
}

// Flow должен быть полный:
// с Connections
// с FlowElements
// с Cursor
// с Workers
func (wf *Workflow) AddFlow(flow *models.Flow) (err error) {

	if flow.Status != "enabled" {
		return
	}

	log.Info("Add flow:", flow.Name)

	wf.mutex.Lock()
	if _, ok := wf.Flows[flow.Id]; ok {
		return
	}
	wf.mutex.Unlock()

	var model *Flow
	if model, err = NewFlow(flow, wf); err != nil {
		return
	}

	wf.mutex.Lock()
	wf.Flows[flow.Id] = model
	wf.mutex.Unlock()


	return
}

func (wf *Workflow) UpdateFlow(flow *models.Flow) (err error) {

	err = wf.RemoveFlow(flow)
	if err != nil {
		return
	}

	err = wf.AddFlow(flow)

	return
}

func (wf *Workflow) RemoveFlow(flow *models.Flow) (err error) {

	log.Info("Remove flow:", flow.Name)

	wf.mutex.Lock()
	defer wf.mutex.Unlock()

	if _, ok := wf.Flows[flow.Id]; !ok {
		return
	}

	wf.Flows[flow.Id].Remove()
	delete(wf.Flows, flow.Id)

	return
}