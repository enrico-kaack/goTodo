package main

import (
	"fmt"
	"github.com/JamesClonk/go-todotxt"
	"log"
)

var TaskList todotxt.TaskList = nil
var ProjectList map[string]int = map[string]int{}
var ContextsList map[string]int = map[string]int{}

func LoadFromFile(fileName string) {
	list, err := todotxt.LoadFromFilename(fileName)
	if err != nil {
		log.Fatal(err)
	}
	TaskList = list
	updateLists()
}

func writeTodoList(filename string) {
	err := todotxt.WriteToFilename(&TaskList, filename)
	if err != nil {
		//show error
	}
}

func updateLists() {
	extractProjectList()
	extractContextList()
}

func extractProjectList() {
	ProjectList = make(map[string]int)
	for _, v := range TaskList {
		for _, p := range v.Projects {
			ProjectList[p] = ProjectList[p] + 1
		}
	}
}

func extractContextList() {
	ContextsList = make(map[string]int)
	for _, v := range TaskList {
		for _, p := range v.Contexts {
			ContextsList[p] = ContextsList[p] + 1
		}
	}
}

func PrintAll() {
	fmt.Print(TaskList)
}

func deleteTaskFromList(index int) {
	err := TaskList.RemoveTask(TaskList[index])
	if err != nil {
		log.Fatal(err)
	}
}

func toggleCompletionState(index int) {
	TaskList[index].Completed = !TaskList[index].Completed
}
