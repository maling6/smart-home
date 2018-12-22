package controllers

import "github.com/gin-gonic/gin"

// simple response
type Response struct {
	StatusCode int
	Payload    interface{}
}

func NewSuccess() *Response {
	resp := &Response{
		StatusCode: 200,
	}
	resp.Success()
	return resp
}

func (r *Response) Success() *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": struct{}{},
	}
	return r
}

func (r *Response) Page(limit, offset, total int, items interface{}) *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": map[string]interface{}{
			"items":  items,
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	}
	return r
}

func (r *Response) Item(name string, item interface{}) *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": map[string]interface{}{
			name: item,
		},
	}
	return r
}

func (r *Response) SetData(data interface{}) *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": data,
	}
	return r
}

func (r *Response) Send(ctx *gin.Context) {
	ctx.JSON(r.StatusCode, r.Payload)
}