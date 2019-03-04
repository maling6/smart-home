package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"errors"
)

func AddWorkflow(workflow *m.Workflow, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, id int64, errs []*validation.Error, err error) {

	// validation
	ok, errs = workflow.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if id, err = adaptors.Workflow.Add(workflow); err != nil {
		return
	}

	workflow.Id = id

	// add workflow
	if workflow, err = adaptors.Workflow.GetById(workflow.Id); err != nil {
		return
	}
	err = core.AddWorkflow(workflow)

	return
}

func GetWorkflowById(workflowId int64, adaptors *adaptors.Adaptors) (workflow *m.Workflow, err error) {

	workflow, err = adaptors.Workflow.GetById(workflowId)

	return
}

func UpdateWorkflow(workflow *m.Workflow, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, errs []*validation.Error, err error) {

	// validation
	ok, errs = workflow.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if err = adaptors.Workflow.Update(workflow); err != nil {
		return
	}

	// reload workflow
	if workflow, err = adaptors.Workflow.GetById(workflow.Id); err != nil {
		return
	}
	err = core.UpdateWorkflow(workflow)

	return
}

func GetWorkflowList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Workflow, total int64, err error) {

	items, total, err = adaptors.Workflow.List(limit, offset, order, sortBy)

	return
}

func DeleteWorkflowById(workflowId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if workflowId == 0 {
		err = errors.New("workflow id is null")
		return
	}

	var workflow *m.Workflow
	if workflow, err = adaptors.Workflow.GetById(workflowId); err != nil {
		return
	}

	// update core
	if err = core.DeleteWorkflow(workflow); err != nil {
		return
	}

	err = adaptors.Workflow.Delete(workflow.Id)

	return
}

func SearchWorkflow(query string, limit, offset int, adaptors *adaptors.Adaptors) (workflows []*m.Workflow, total int64, err error) {

	workflows, total, err = adaptors.Workflow.Search(query, limit, offset)

	return
}

func UpdateWorkflowScenario(workflow *m.Workflow, workflowScenarioId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if err = adaptors.Workflow.SetScenario(workflow, workflowScenarioId); err != nil {
		return
	}

	// update core
	if workflow, err = adaptors.Workflow.GetById(workflow.Id); err != nil {
		return
	}
	core.UpdateWorkflowScenario(workflow)

	return
}
