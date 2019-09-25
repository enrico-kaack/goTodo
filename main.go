package main

import (
	"fmt"
	"github.com/JamesClonk/go-todotxt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
)

var filePath string = "todo.txt"

var horizontalFlex = tview.NewFlex().SetDirection(tview.FlexRow)
var verticalFlex = tview.NewFlex()
var leftList = tview.NewFlex().SetDirection(tview.FlexRow)
var pages = tview.NewPages()

var taskListView = tview.NewList().ShowSecondaryText(false)
var contextListView = tview.NewList().ShowSecondaryText(false)
var projectListView = tview.NewList().ShowSecondaryText(false)

var addInputView = tview.NewInputField()
var editInputView = tview.NewInputField()

var priorityInputView = tview.NewInputField()

var currentPage = "main"

var app = tview.NewApplication()

func main() {
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	LoadFromFile(filePath)
	renderLists()

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 3, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	contextListView.SetBorder(true).SetTitle("Contexts")
	projectListView.SetBorder(true).SetTitle("Projects")
	taskListView.SetBorder(true).SetTitle("Tasks").SetInputCapture(taskListKeyHandler)

	addInputView.SetLabel("Add: ")
	addInputView.SetDoneFunc(addInputDoneHandler)
	editInputView.SetLabel("Edit: ")
	editInputView.SetDoneFunc(editInputDoneHandler)
	priorityInputView.SetDoneFunc(priorityInputDoneHandler).SetBorder(true).SetTitle("Priority")

	leftList.AddItem(contextListView, 0, 1, false).AddItem(projectListView, 0, 1, false)
	verticalFlex.AddItem(leftList, 20, 1, false)
	verticalFlex.AddItem(taskListView, 0, 5, false)

	horizontalFlex.AddItem(verticalFlex, 0, 1, false)
	pages.AddPage("main", horizontalFlex, true, true)
	pages.AddPage("editPriority", modal(priorityInputView, 20, 3), true, false)
	pages.AddPage("detail", modal(createDetailPage(), 70, 0), true, false)

	app.SetInputCapture(keyHandler)
	if err := app.SetRoot(pages, true).SetFocus(taskListView).Run(); err != nil {
		panic(err)
	}
}

func renderLists() {
	updateLists()

	taskListView.Clear()
	for _, v := range TaskList {
		taskListView.AddItem(fmt.Sprintf("%s %s %s", FormatDone(v.Completed), v.Priority, v.Todo), "", 0, nil)
	}

	projectListView.Clear()
	for projectName, projectCount := range ProjectList {
		projectListView.AddItem(fmt.Sprintf("%s [%d]", projectName, projectCount), "", 0, nil)
	}

	contextListView.Clear()
	for contextName, contextCount := range ContextsList {
		contextListView.AddItem(fmt.Sprintf("%s [%d]", contextName, contextCount), "", 0, nil)
	}
}

func taskListKeyHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEnter {
		pages.SwitchToPage("detail")
		currentPage = "detail"

		setDetailPageData(TaskList[taskListView.GetCurrentItem()])
	}

	switch event.Rune() {
	case 'x':
		index := taskListView.GetCurrentItem()
		deleteTaskFromList(index)
		renderLists()
		taskListView.SetCurrentItem(index - 1)
	case 'd':
		index := taskListView.GetCurrentItem()
		toggleCompletionState(index)
		renderLists()
		taskListView.SetCurrentItem(index)
	case 'e':
		horizontalFlex.AddItem(editInputView, 1, 0, false)
		editInputView.SetText(TaskList[taskListView.GetCurrentItem()].String())
		app.SetFocus(editInputView)
	case 'p':
		pages.SwitchToPage("editPriority")
		currentPage = "editPriority"
		app.SetFocus(priorityInputView)

	}

	return event
}

func keyHandler(event *tcell.EventKey) *tcell.EventKey {
	if currentPage != "main" && currentPage != "editPriority" {
		return event
	}
	if app.GetFocus() == addInputView || app.GetFocus() == editInputView || app.GetFocus() == priorityInputView {
		return event
	}

	switch event.Rune() {

	case 'T':
		app.SetFocus(taskListView)
	case 'C':
		app.SetFocus(contextListView)
	case 'P':
		app.SetFocus(projectListView)
	case 'a':
		horizontalFlex.AddItem(addInputView, 1, 0, false)
		app.SetFocus(addInputView)
	case 's':
		writeTodoList(filePath)

	}

	return event
}

func priorityInputDoneHandler(key tcell.Key) {
	if key == tcell.KeyEsc {
		resetAndHidePriorityInputField()
		return
	}

	index := taskListView.GetCurrentItem()
	TaskList[index].Priority = priorityInputView.GetText()
	renderLists()
	resetAndHidePriorityInputField()
}

func editInputDoneHandler(key tcell.Key) {
	if key == tcell.KeyEsc {
		resetAndHideEditInputField()
		return
	}

	editedTask, err := todotxt.ParseTask(editInputView.GetText())
	if err != nil {
		//TODO: error handling
		resetAndHideEditInputField()
		return
	}
	index := taskListView.GetCurrentItem()
	TaskList[index] = *editedTask
	renderLists()
	taskListView.SetCurrentItem(index)
	resetAndHideEditInputField()

}

func addInputDoneHandler(key tcell.Key) {
	if key == tcell.KeyEsc {
		resetAndHideAddInputField()
		return
	}

	newTask, err := todotxt.ParseTask(addInputView.GetText())
	if err != nil {
		//handle format problem
	}
	TaskList.AddTask(newTask)
	renderLists()

	resetAndHideAddInputField()
}

func resetAndHidePriorityInputField() {
	priorityInputView.SetText("")
	pages.SwitchToPage("main")
	currentPage = "main"
	app.SetFocus(taskListView)
}

func resetAndHideEditInputField() {
	horizontalFlex.RemoveItem(editInputView)
	editInputView.SetText("")
	app.SetFocus(taskListView)
}

func resetAndHideAddInputField() {
	horizontalFlex.RemoveItem(addInputView)
	addInputView.SetText("")
	app.SetFocus(taskListView)
}
