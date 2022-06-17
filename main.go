package main

import "github.com/moson-mo/reminder/internal/reminder"

func main() {
	err := reminder.Start()
	if err != nil {
		panic(err)
	}
}
