package cmd

import (
	"fmt"
	"os"

	"github.com/disosur/gid/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Switch to a GitHub profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]

		cfg, err := internal.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, color.RedString("Error: %v", err))
			os.Exit(1)
		}

		profile, ok := cfg.Profiles[alias]
		if !ok {
			fmt.Fprintln(os.Stderr, color.RedString("Profile %q not found", alias))
			os.Exit(1)
		}

		if err := internal.SetGitUser(profile.Name, profile.Email); err != nil {
			fmt.Fprintln(os.Stderr, color.RedString("Failed to set git user: %v", err))
			os.Exit(1)
		}

		if err := internal.SetGitSSH(profile.SSHKey); err != nil {
			fmt.Fprintln(os.Stderr, color.RedString("Failed to set SSH key: %v", err))
			os.Exit(1)
		}

		cfg.Active = alias
		if err := internal.SaveConfig(cfg); err != nil {
			fmt.Fprintln(os.Stderr, color.RedString("Failed to save config: %v", err))
			os.Exit(1)
		}

		color.Green("✔ Switched to %s (%s <%s>)", alias, profile.Name, profile.Email)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
