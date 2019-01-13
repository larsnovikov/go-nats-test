package main

import (
	"./helpers"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func main() {
	go startHandlers()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		route(writer, request)
	})

	log.Fatal(http.ListenAndServe(":9107", nil))
}

func route(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	Helper.ErrorHandler.Handle(err, func(err error) {}, func(err error) {})

	if string(request.URL.Path[len(request.URL.Path)-1:]) != "/" {
		http.Redirect(writer, request, request.URL.Path+"/", 301)
	} else {
		var routeName string

		for key := range Helper.Routing {
			match, _ := regexp.MatchString(key+"$", request.URL.Path)
			if match {
				routeName = key
				break
			}
		}

		if routeName != "" {
			var result map[string]interface{}
			var res []byte

			if Helper.Routing[routeName].AfterPublish != nil {
				// Выполнение через NATS (изменение и удаление данных)
				Helper.NatsHandler.Publish(Helper.Routing[routeName])
				result = Helper.Routing[routeName].AfterPublish(Helper.Routing[routeName].Request)
			} else {
				// Выполнение напрямую
				result = Helper.Routing[routeName].Handler(Helper.Routing[routeName].Request)
			}

			res, err = json.Marshal(result)
			Helper.ErrorHandler.Handle(err, func(err error) {}, func(err error) {})

			writer.Header().Set("Content-Type", "application/json")
			writer.Write(res)
		}
	}
}

func startHandlers() {
	for i, handler := range Helper.Config {
		handler.ClientID = "client" + strconv.Itoa(i)
		handler.Subscribe()
	}
}
