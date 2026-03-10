package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test <profile>",
	Short: "Test SSH connection to GitHub for a profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		host := fmt.Sprintf("git@github.com-%s", alias)

		color.Cyan("Testing SSH connection for %q → %s ...", alias, host)

		c := exec.Command("ssh", "-T", host)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		// ssh -T git@github.com exits 1 even on success, so we just
		// print the output and let the user interpret it.
		_ = c.Run()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
