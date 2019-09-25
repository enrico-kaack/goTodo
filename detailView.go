package main

import (
	"github.com/JamesClonk/go-todotxt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var detailView = tview.NewTextView()
var priorityView = tview.NewTextView()
var contextView = tview.NewList()
var projectView = tview.NewList()
var doneView = tview.NewTextView()
var startDateView = tview.NewTextView()
var dueDateView = tview.NewTextView()

var Task todotxt.Task

func createDetailPage() tview.Primitive {
	priorityView.SetBorder(true).SetTitle("P")
	startDateView.SetBorder(true).SetTitle("Start")
	dueDateView.SetBorder(true).SetTitle("Due")
	doneView.SetBorder(true).SetTitle("Status")
	detailView.SetBorder(true).SetTitle("ToDo")
	detailView.SetScrollable(true)
	projectView.SetBorder(true).SetTitle("Projects")
	projectView.SetSelectedFocusOnly(true)
	contextView.SetBorder(true).SetTitle("Contexts")
	contextView.SetSelectedFocusOnly(true)

	horizontalFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	firstFlex := tview.NewFlex().AddItem(priorityView, 0, 1, false)
	firstFlex.AddItem(doneView, 0, 2, false)
	firstFlex.AddItem(startDateView, 0, 3, false)
	firstFlex.AddItem(dueDateView, 0, 3, false)
	horizontalFlex.AddItem(firstFlex, 3, 1, false)
	horizontalFlex.AddItem(detailView, 4, 3, true)

	secondFlex := tview.NewFlex().AddItem(projectView, 0, 1, false)
	secondFlex.AddItem(contextView, 0, 1, false)
	horizontalFlex.AddItem(secondFlex, 0, 3, false)
	detailView.SetInputCapture(detailPageInputHandler)
	return horizontalFlex
}

func setDetailPageData(task todotxt.Task) {
	priorityView.SetText(task.Priority)
	doneView.SetText(FormatDoneInString(task.Completed))
	startDateView.SetText(FormatDate(task.CreatedDate))
	dueDateView.SetText(FormatDate(task.DueDate))
	detailView.SetText(task.Todo)

	projectView.Clear()
	for _, v := range task.Projects {
		projectView.AddItem(v, "", 0, nil)
	}
	contextView.Clear()
	for _, v := range task.Contexts {
		contextView.AddItem(v, "", 0, nil)
	}

	app.SetFocus(detailView)
}

func detailPageInputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		pages.SwitchToPage("main")
		currentPage = "main"
		app.SetFocus(taskListView)
	}

	return event
}
