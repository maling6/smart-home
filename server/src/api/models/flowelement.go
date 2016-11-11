package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"sync"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type FlowElement struct {
	Id   		int64  			`orm:"pk;auto;column(id)" json:"id"`
	Name		string			`orm:"" json:"name"`
	GraphSettings	string			`orm:"column(graph_settings)" json:"graph_settings"`
	Status		string			`orm:"" json:"status"`
	FlowId		int64			`orm:"column(flow_id)" json:"flow_id"`
	Created_at	time.Time		`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time		`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
	PrototypeType	string			`orm:"column(prototype_type)" json:"prototype_type"`
	Flow		*Flow			`orm:"-" json:"-"`
	Prototype	ActionPrototypes	`orm:"-" json:"-"`
	status		Status			`orm:"-" json:"-"`
	mutex     	sync.Mutex		`orm:"-" json:"-"`
}

func (m *FlowElement) TableName() string {
	return beego.AppConfig.String("db_flow_elements")
}

func init() {
	orm.RegisterModel(new(FlowElement))
}

// AddFlowElement insert a new FlowElement into database and returns
// last inserted Id on success.
func AddFlowElement(m *FlowElement) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetFlowElementById retrieves FlowElement by Id. Returns error if
// Id doesn't exist
func GetFlowElementById(id int64) (v *FlowElement, err error) {
	o := orm.NewOrm()
	v = &FlowElement{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllFlowElement retrieves all FlowElement matches certain condition. Returns empty list if
// no records exist
func GetAllFlowElement(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(FlowElement))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []FlowElement
	qs = qs.OrderBy(sortFields...)
	objects_count, err := qs.Count()
	if err != nil {
		return
	}
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		meta = &map[string]int64{
			"objects_count": objects_count,
			"limit":         limit,
			"offset":        offset,
		}
		return ml, meta, nil
	}
	return nil, nil, err
}

// UpdateFlowElement updates FlowElement by Id and returns error if
// the record to be updated doesn't exist
func UpdateFlowElementById(m *FlowElement) (err error) {
	o := orm.NewOrm()
	v := FlowElement{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFlowElement deletes FlowElement by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFlowElement(id int64) (err error) {
	o := orm.NewOrm()
	v := FlowElement{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&FlowElement{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetFlowElementsByFlow(flow *Flow) (elements []*FlowElement, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(&FlowElement{}).Filter("status", "enabled").Filter("flow_id", flow.Id).All(&elements)
	return
}

//---------------------------------------------------
// Workflow
//---------------------------------------------------
func (m *FlowElement) Before(message *Message) error {

	if m.Flow == nil {
		return errors.New("flow is nil ...")
	}
	// mutex
	//----------------------
	m.mutex.Lock()
	m.Flow.Cursor  = append(m.Flow.Cursor, m)
	m.mutex.Unlock()

	m.status = DONE
	return  m.Prototype.Before(message, m.Flow)
}

// run internal process
func (m *FlowElement) Run(message *Message) (err error) {
	m.status = IN_PROCESS

	if m.Flow == nil {
		return errors.New("flow is nil ...")
	}

	err = m.Before(message)
	err = m.Prototype.Run(message, m.Flow)
	err = m.After(message)

	if err != nil {
		return
	}

	var elements []*FlowElement
	for _, conn := range m.Flow.Connections {
		if conn.ElementFrom != m.Id {
			continue
		}

		elements = append(elements, conn.FlowElementTo)
	}

	for _, element := range elements {
		if len(elements) > 1 {
			go element.Run(message)
		} else {
			element.Run(message)
		}
	}

	m.status = ENDED

	if len(elements) == 0 {
		err = errors.New("Element not found")
	}

	return
}

func (m *FlowElement) After(message *Message) error {
	m.status = STARTED

	if m.Flow == nil {
		return errors.New("flow is nil ...")
	}

	// mutex
	//----------------------
	m.mutex.Lock()
	if len(m.Flow.Cursor) > 1 {
		for k, f := range m.Flow.Cursor {
			if f != m {
				continue
			}

			m.Flow.Cursor = append(m.Flow.Cursor[:k], m.Flow.Cursor[k+1:]...)
		}
	} else {
		m.Flow.Cursor = []*FlowElement{}
	}
	m.mutex.Unlock()

	return  m.Prototype.After(message, m.Flow)
}
