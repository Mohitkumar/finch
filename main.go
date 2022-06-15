package main

import (
	"fmt"

	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence/standalone"
)

func main() {
	conf := &standalone.Config{
		DBPath: "/Users/mohitkumar/bdg",
	}
	dao, err := standalone.NewBadgerWorkflowDaoImpl(*conf)
	if err != nil {
		panic(err)
	}

	wf := &model.Workflow{
		Name:       "wf1",
		RootAction: 1,
	}
	status, _ := dao.Save(*wf)
	if status {
		dd, _ := dao.Get("wf1")
		fmt.Print(dd)
	}
}
