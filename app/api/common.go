package api

//go:generate easyjson

import (
	"encoding/json"
	"fmt"

	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/goppy/plugins/web"
	"github.com/osspkg/hermes/app/acl"
	"github.com/osspkg/hermes/app/addons"
	"github.com/osspkg/hermes/app/collections"
)

type API struct {
	addons      *addons.Addons
	collections *collections.Collections
	acl         *acl.ACL
	router      web.Router
}

func New(a *addons.Addons, c *collections.Collections, ac *acl.ACL, r web.RouterPool) *API {
	return &API{
		addons:      a,
		collections: c,
		acl:         ac,
		router:      r.Main(),
	}
}

func (v *API) Up(ctx app.Context) error {
	v.router.Post("/api/jsonrpc", v.JsonRPC)
	return nil
}

func (v *API) Down() error {
	return nil
}

type (
	//easyjson:json
	RequestModel struct {
		Addon string `json:"addon"`
		Type  uint   `json:"type"`
		Form  uint   `json:"form"`
		Data  []byte `json:"data"`
	}
	//easyjson:json
	ResponseModel struct {
		Data  []byte `json:"data"`
		Error string `json:"error"`
	}
)

func (v *API) JsonRPC(ctx web.Context) {
	model := &RequestModel{}
	if err := ctx.BindJSON(model); err != nil {
		response(ctx, nil, err)
		return
	}
	api, err := v.addons.ResolveApi(model.Addon)
	if err != nil {
		response(ctx, nil, err)
		return
	}
	var data json.Marshaler
	if model.Type == 0 {
		data, err = api.Form(model.Form)
	} else if model.Type == 1 {
		data, err = api.Call(ctx.Context(), model.Form, model.Data, nil)
	} else {
		response(ctx, nil, fmt.Errorf("unknown request type"))
		return
	}
	if err != nil {
		response(ctx, nil, err)
		return
	}
	b, er := data.MarshalJSON()
	response(ctx, b, er)
}

func response(ctx web.Context, d []byte, err error) {
	model := &ResponseModel{}
	if err != nil {
		model.Error = err.Error()
		ctx.JSON(400, model)
		return
	}
	model.Data = d
	ctx.JSON(200, model)
}
