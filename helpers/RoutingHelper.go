package Helper

import (
	"../controllers"
	"../structs"
)

var Routing = map[string]structs.Route{
	"/user/": {
		Topic:      "get_data",
		Controller: "user",
		Action:     "list",
		Handler:    controllers.UserController.Detail,
		Request: structs.Request{
			Method:      "GET",
			Path:        "/user/list",
			PathPattern: "/user/",
		},
	},
	"/user/([a-z,0-9]+)/": {
		Topic:      "get_data",
		Controller: "user",
		Action:     "detail",
		Handler:    controllers.UserController.Detail,
		Request: structs.Request{
			Method:      "GET",
			Path:        "/user/detail/",
			PathPattern: "/user/([a-z,0-9]+)/",
		},
	},
	"/user/delete([a-z,0-9]+)/": {
		Topic:      "set_data",
		Controller: "user",
		Action:     "delete",
		Handler:    controllers.UserController.Delete,
		Request: structs.Request{
			Method:      "POST",
			Path:        "/user/delete/",
			PathPattern: "/user/delete/([a-z,0-9]+)/",
		},
		AfterPublish: func(request structs.Request) map[string]interface{} {
			res := make(map[string]interface{})
			res["name"] = "User delete published"

			return res
		},
	},
	"/user/update([a-z,0-9]+)/": {
		Topic:      "set_data",
		Controller: "user",
		Action:     "update",
		Handler:    controllers.UserController.Update,
		Request: structs.Request{
			Method:      "POST",
			Path:        "/user/update/",
			PathPattern: "/user/update([a-z,0-9]+)/",
		},
		AfterPublish: func(request structs.Request) map[string]interface{} {
			res := make(map[string]interface{})
			res["name"] = "User update published"

			return res
		},
	},
	"/post/": {
		Topic:      "get_data",
		Controller: "post",
		Action:     "list",
		Handler:    controllers.PostController.List,
		Request: structs.Request{
			Method:      "GET",
			Path:        "/post/list/",
			PathPattern: "/post/",
		},
	},
	"/post/([a-z,0-9]+)/": {
		Topic:      "get_data",
		Controller: "post",
		Action:     "detail",
		Handler:    controllers.PostController.Detail,
		Request: structs.Request{
			Method:      "GET",
			Path:        "/post/detail/",
			PathPattern: "/post/([a-z,0-9]+)/",
		},
	},
}
