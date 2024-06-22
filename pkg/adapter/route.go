package adapter

import (
	"fmt"

	"github.com/DictumMortuum/servus-extapi/pkg/middleware"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
)

func RaRoute(router *gin.RouterGroup, endpoint string, obj model.Routable, funcs ...func(*gin.Context)) {
	group := router.Group("/" + endpoint)
	group.Use(func(c *gin.Context) {
		m, err := model.ToMap(c, "req")
		if err != nil {
			c.Error(err)
			return
		}

		m.Set("apimodel", obj)
	})

	// jwt := middleware.Jwt("http://sol.dictummortuum.com:3567/.well-known/jwks.json")

	listgroup := group.Group("")
	for _, fn := range funcs {
		listgroup.Use(fn)
	}

	{
		listgroup.GET(
			"",
			middleware.Filter,
			middleware.Paginate,
			middleware.Sort,
			A(CountMany),
			A(GetMany),
			middleware.ResultRa,
		)
		group.GET(
			"/:id",
			middleware.Id,
			A(GetOne),
			middleware.ResultRa,
		)
		group.PUT(
			"/:id",
			// jwt,
			middleware.Id,
			middleware.Body,
			A(UpdateOne),
			middleware.ResultRa,
		)
		group.POST(
			"",
			// jwt,
			middleware.Body,
			A(CreateOne),
			middleware.ResultRa,
		)
		group.DELETE(
			"/:id",
			// jwt,
			middleware.Id,
			A(DeleteOne),
			middleware.ResultRa,
		)
	}
}

func Route(router *gin.RouterGroup, endpoint string, obj model.Routable) {
	group := router.Group("/" + endpoint)
	group.Use(func(c *gin.Context) {
		m, err := model.ToMap(c, "req")
		if err != nil {
			c.Error(err)
			return
		}

		m.Set("apimodel", obj)
	})
	{
		group.GET(
			"",
			middleware.Filter,
			middleware.Paginate,
			middleware.Sort,
			A(CountMany),
			A(GetMany),
			middleware.Result,
		)
		group.GET(
			"/:id",
			middleware.Id,
			A(GetOne),
			middleware.Result,
		)
		group.PUT(
			"/:id",
			middleware.Id,
			middleware.Body,
			A(UpdateOne),
			middleware.Result,
		)
		group.POST(
			"",
			middleware.Body,
			A(CreateOne),
			middleware.Result,
		)
		group.DELETE(
			"/:id",
			middleware.Id,
			A(DeleteOne),
			middleware.Result,
		)
	}
}

func GetOne(req *model.Map, res *model.Map) error {
	m, err := req.GetModel()
	if err != nil {
		return err
	}

	DB, err := req.GetGorm()
	if err != nil {
		return err
	}

	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	data, err := m.Get(DB, id)
	if err != nil {
		return err
	}

	res.Set("data", data)
	return nil
}

func CountMany(req *model.Map, res *model.Map) error {
	m, err := req.GetModel()
	if err != nil {
		return err
	}

	limits, err := req.GetString("range")
	if err != nil {
		return err
	}

	DB, err := req.GetGorm()
	if err != nil {
		return err
	}

	var count int64
	rs := DB.Model(m).Scopes(m.DefaultFilter, req.Filter).Count(&count)
	if rs.Error != nil {
		return rs.Error
	}
	res.Headers["Content-Range"] = fmt.Sprintf("%s/%d", limits, count)

	return nil
}

func GetMany(req *model.Map, res *model.Map) error {
	m, err := req.GetModel()
	if err != nil {
		return err
	}

	DB, err := req.GetGorm()
	if err != nil {
		return err
	}

	data, err := m.List(DB, m.DefaultFilter, req.Filter, req.Sort, req.Paginate)
	if err != nil {
		return err
	}

	res.Set("data", data)
	return nil
}

func UpdateOne(req *model.Map, res *model.Map) error {
	m, err := req.GetModel()
	if err != nil {
		return err
	}

	DB, err := req.GetGorm()
	if err != nil {
		return err
	}

	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	body, err := req.GetByte("body")
	if err != nil {
		return err
	}

	data, err := m.Update(DB, id, body)
	if err != nil {
		return err
	}

	res.Set("data", data)
	return nil
}

func CreateOne(req *model.Map, res *model.Map) error {
	m, err := req.GetModel()
	if err != nil {
		return err
	}

	DB, err := req.GetGorm()
	if err != nil {
		return err
	}

	body, err := req.GetByte("body")
	if err != nil {
		return err
	}

	data, err := m.Create(DB, body)
	if err != nil {
		return err
	}

	res.Set("data", data)
	return nil
}

func DeleteOne(req *model.Map, res *model.Map) error {
	m, err := req.GetModel()
	if err != nil {
		return err
	}

	DB, err := req.GetGorm()
	if err != nil {
		return err
	}

	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	data, err := m.Delete(DB, id)
	if err != nil {
		return err
	}

	res.Set("data", data)
	return nil
}
