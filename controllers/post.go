package controllers

import (
	"../structs"
)

type postController struct {
}

var PostController postController

func (p postController) Detail(request structs.Request) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = "Detail POST"

	return res
}

func (p postController) List(request structs.Request) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = "List POST"
	return res
}
