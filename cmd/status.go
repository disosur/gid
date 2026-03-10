package cmd

import (
	"fmt"
	"os"

	"github.com/disosur/gid/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the currently active GitHub profile",
	Run: func(cmd *cobra.Command, args []string) {
		email, err := internal.GetGitEmail()
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString("Could not read git email: %v", err))
			os.Exit(1)
		}

		name, _ := internal.GetGitName()

		cfg, err := internal.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString("Error: %v", err))
			os.Exit(1)
		}

		for alias, p := range cfg.Profiles {
			if p.Email == email {
				color.Green("Active: %s (%s <%s>)", alias, name, email)
				return
			}
		}

		color.Yellow("No matching profile for %s <%s>", name, email)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
