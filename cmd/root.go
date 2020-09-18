package cmd

import (
	"fmt"
	"os"

	"github.com/1995parham/zamaneh/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// nolint: gocheckglobals
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	// ExitFailure status code.
	ExitFailure = 1

	// TopicFlag is the name of topic flag.
	TopicFlag = "topic"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	root := &cobra.Command{
		Use:     "zamaneh",
		Short:   "Manage your working periods with ease",
		Example: "zamaneh --topic ml",
		Version: fmt.Sprintf("%s %s [%s]", version, commit, date),
		Args:    cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			topic, err := cmd.Flags().GetString(TopicFlag)
			if err != nil {
				panic(err)
			}

			main(topic)
		},
	}

	root.Flags().StringP(TopicFlag, "t", "untitled", "working period title")

	if err := root.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(ExitFailure)
	}
}

func main(topic string) {
	a, err := app.New(topic)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := a.Run(); err != nil {
		logrus.Fatal(err)
	}
}
