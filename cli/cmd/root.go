package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "fil-vote",
	Short:   "fil-vote cli tool",
	Long:    "fil-vote cli tool for managing voting-related tasks",
	Version: "1.0.0",
}
