package Helper

import "fmt"

type errorHelper struct {
}

var ErrorHandler errorHelper

func (e errorHelper) Handle(err error, before func(err error), after func(err error)) {
	if err != nil {
		before(err)
		fmt.Println(err)
		after(err)
	}
}
