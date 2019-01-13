package controllers

import (
	"../structs"
)

type userController struct {
}

var UserController userController

func (u userController) Detail(request structs.Request) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = "Detail USER"
	res["content"] = request.Path

	return res
}

func (u userController) List(request structs.Request) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = "List USER"
	res["content"] = request.Path
	return res
}

func (u userController) Delete(request structs.Request) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = "Delete USER"
	res["content"] = request.Path
	return res
}

func (u userController) Update(request structs.Request) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = "Update USER"
	res["content"] = request.Path
	return res
}
