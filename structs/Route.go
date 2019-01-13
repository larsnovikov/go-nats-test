package structs

type Route struct {
	Topic        string `json:"topic"`
	Controller   string `json:"controller"`
	Action       string `json:"action"`
	Handler      func(request Request) map[string]interface{}
	AfterPublish func(request Request) map[string]interface{}
	Request
}
