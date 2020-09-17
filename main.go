package main

import (
	"os"

	"github.com/1995parham/zamaneh/app"
	"github.com/sirupsen/logrus"
)

func main() {
	topic := "untitled"
	if len(os.Args) > 1 {
		topic = os.Args[1]
	}

	a, err := app.New(topic)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := a.Run(); err != nil {
		logrus.Fatal(err)
	}
}
