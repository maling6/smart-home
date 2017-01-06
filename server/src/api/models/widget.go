package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
)

type Widget struct {
	Id           	int64  			`orm:"pk;auto;column(id)" json:"id"`
	Name        	string			`orm:"" json:"name" valid:"MaxSize(254);Required"`
	Dashboard	*Dashboard		`orm:"rel(fk)" json:"dashboard"`
	Width		int8			`orm:"" json:"width"`
	Height		int8			`orm:"" json:"height"`
	Column		int8			`orm:"" json:"column"`
	Row		int8			`orm:"" json:"row"`
	Created_at   	time.Time		`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at    	time.Time		`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *Widget) TableName() string {
	return beego.AppConfig.String("db_widgets")
}

func init() {
	orm.RegisterModel(new(Widget))
}

// AddWidget insert a new Widget into database and returns
// last inserted Id on success.
func AddWidget(m *Widget) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetWidgetById retrieves Widget by Id. Returns error if
// Id doesn't exist
func GetWidgetById(id int64) (v *Widget, err error) {
	o := orm.NewOrm()
	v = &Widget{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllWidget retrieves all Widget matches certain condition. Returns empty list if
// no records exist
func GetAllWidget(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Widget))
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

	var l []Widget
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

// UpdateWidget updates Widget by Id and returns error if
// the record to be updated doesn't exist
func UpdateWidgetById(m *Widget) (err error) {
	o := orm.NewOrm()
	v := Widget{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteWidget deletes Widget by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWidget(id int64) (err error) {
	o := orm.NewOrm()
	v := Widget{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Widget{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
