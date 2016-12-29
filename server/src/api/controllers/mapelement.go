package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"fmt"
	"../models"
)

// MapElementController operations for MapElement
type MapElementController struct {
	CommonController
}

// URLMapping ...
func (c *MapElementController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("PutElementOnly", c.PutElementOnly)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create MapElement
// @Param	body		body 	models.MapElement	true		"body for MapElement content"
// @Success 201 {object} models.MapElement
// @Failure 403 body is empty
// @router / [post]
func (c *MapElementController) Post() {
	var element models.MapElement
	json.Unmarshal(c.Ctx.Input.RequestBody, &element)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&element)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if !b {
		var msg string
		for _, err := range valid.Errors {
			msg += fmt.Sprintf("%s: %s\r", err.Key, err.Message)
		}
		c.ErrHan(403, msg)
		return
	}

	switch element.PrototypeType {
	case "text":
	case "image":
		prototype := &models.MapImage{}
		b, _ := json.Marshal(element.Prototype)
		json.Unmarshal(b, &prototype)
		if prototype != nil {
			if prototype.Id == 0 {
				if _, err := models.AddMapImage(prototype); err != nil {
					c.ErrHan(403, err.Error())
					return
				}

				element.PrototypeId = prototype.Id
				element.PrototypeType = "image"
			}
		}
	case "device":
	case "script":
	default:

	}

	if _, err = models.AddMapElement(&element); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map_element": element}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get MapElement by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.MapElement
// @Failure 403 :id is empty
// @router /:id [get]
func (c *MapElementController) GetOne() {
	id, _ := c.GetInt(":id")
	entity, err := models.GetMapElementById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"map_element": entity}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get MapElement
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.MapElement
// @Failure 403
// @router / [get]
func (c *MapElementController) GetAll() {
	ml, meta, err := models.GetAllMap(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"map_entities": ml, "meta": meta}
	c.ServeJSON()
}

// PutElementOnly ...
// @Title PutElementOnly
// @Description update the MapElement
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.MapElement	true		"body for MapElement content"
// @Success 200 {object} models.MapElement
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MapElementController) PutElementOnly() {

	var mapElement, oldMapElement *models.MapElement
	var err error

	id, _ := c.GetInt(":id")

	if mapElement, err = models.GetMapElementById(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}
	oldMapElement = &models.MapElement{}
	*oldMapElement = *mapElement

	json.Unmarshal(c.Ctx.Input.RequestBody, &mapElement)
	mapElement.Id = int64(id)
	mapElement.PrototypeType = oldMapElement.PrototypeType
	mapElement.PrototypeId = oldMapElement.PrototypeId

	// update map element
	//
	if err := models.UpdateMapElementById(mapElement); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the MapElement
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.MapElement	true		"body for MapElement content"
// @Success 200 {object} models.MapElement
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MapElementController) Put() {
	id, _ := c.GetInt(":id")
	var mapElement models.MapElement
	json.Unmarshal(c.Ctx.Input.RequestBody, &mapElement)
	mapElement.Id = int64(id)

	// load old data
	//
	var oldMapElement *models.MapElement
	var err error
	if oldMapElement, err = models.GetMapElementById(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	// delete old map image
	//
	if oldMapElement.PrototypeId == 0 {
		oldMapElement.PrototypeType = ""
	}

	switch oldMapElement.PrototypeType {
	case "text":
	case "image":
		models.DeleteMapImage(oldMapElement.PrototypeId)
	case "device":
	case "script":
	}

	//
	switch mapElement.PrototypeType {
	case "text":
	case "image":

		prototype := &models.MapImage{}
		b, _ := json.Marshal(mapElement.Prototype)
		json.Unmarshal(b, &prototype)
		if prototype != nil {
			if prototype.Id == 0 {
				if _, err := models.AddMapImage(prototype); err != nil {
					c.ErrHan(403, err.Error())
					return
				}
			} else {
				if err := models.UpdateMapImageById(prototype); err != nil {
					c.ErrHan(403, err.Error())
					return
				}
			}

			mapElement.PrototypeId = prototype.Id
			mapElement.PrototypeType = "image"
		} else {
			mapElement.PrototypeId = 0
			mapElement.PrototypeType = ""
		}
	case "device":
	case "script":

	}

	// update map element
	//
	if err := models.UpdateMapElementById(&mapElement); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the MapElement
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *MapElementController) Delete() {
	id, _ := c.GetInt(":id")

	var err error
	var oldMapElement *models.MapElement
	if oldMapElement, err = models.GetMapElementById(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	switch oldMapElement.PrototypeType {
	case "image":
		models.DeleteMapImage(oldMapElement.PrototypeId)
	case "device":
	case "script":
	}

	if err := models.DeleteMapElement(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}
