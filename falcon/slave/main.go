package main

import (
	"falcon/slave/cases"
	"falcon/slave/site/boomer"
	"os"
)

func main() {
	boomer.LoadArgv()
	boomer.CaseEntry = cases.NewCase
	boomer.RunTasks()
	os.Exit(0)
}
