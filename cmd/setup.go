package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/disosur/gid/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive wizard to configure GitHub profiles",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		count := promptInt(reader, "How many profiles do you want to configure?")
		if count <= 0 {
			fmt.Println("Nothing to configure.")
			return
		}

		cfg := &internal.Config{
			Profiles: make(map[string]internal.Profile),
		}

		for i := 1; i <= count; i++ {
			fmt.Printf("\n--- Profile %d ---\n", i)

			alias := prompt(reader, "Profile alias (e.g. tulay, work)")
			name := prompt(reader, "Full name")
			email := prompt(reader, "Email")

			genKey := promptYN(reader, "Generate SSH key?", true)

			var keyPath string
			if genKey {
				path, err := internal.GenerateSSHKey(alias, email)
				if err != nil {
					color.Yellow("⚠ %v", err)
					if path != "" {
						keyPath = path // key already exists, reuse it
					} else {
						keyPath = prompt(reader, "SSH key path (e.g. ~/.ssh/id_ed25519)")
					}
				} else {
					keyPath = path
					color.Green("✔ Generated %s", keyPath)

					if err := internal.AddSSHHostBlock(alias, keyPath); err != nil {
						color.Red("Failed to add SSH config block: %v", err)
					} else {
						color.Green("✔ Added Host block to ~/.ssh/config")
					}

					pubKey, err := internal.ReadPublicKey(keyPath)
					if err == nil {
						fmt.Println()
						color.Cyan("Add this public key to github.com/settings/keys:")
						fmt.Println()
						fmt.Println(pubKey)
						fmt.Println()
						prompt(reader, "Press Enter when done...")
					}
				}
			} else {
				keyPath = prompt(reader, "SSH key path (e.g. ~/.ssh/id_ed25519)")
			}

			cfg.Profiles[alias] = internal.Profile{
				Name:   name,
				Email:  email,
				SSHKey: keyPath,
			}

			// first profile becomes active by default
			if cfg.Active == "" {
				cfg.Active = alias
			}
		}

		if err := internal.SaveConfig(cfg); err != nil {
			color.Red("Failed to save config: %v", err)
			os.Exit(1)
		}

		color.Green("\n✔ Setup complete! %d profile(s) saved.", count)
		color.Cyan("Run `gid use <profile>` to switch.")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func prompt(r *bufio.Reader, label string) string {
	fmt.Printf("%s: ", label)
	text, _ := r.ReadString('\n')
	return strings.TrimSpace(text)
}

func promptInt(r *bufio.Reader, label string) int {
	s := prompt(r, label)
	n, _ := strconv.Atoi(s)
	return n
}

func promptYN(r *bufio.Reader, label string, defaultYes bool) bool {
	hint := "[Y/n]"
	if !defaultYes {
		hint = "[y/N]"
	}
	fmt.Printf("%s %s: ", label, hint)
	text, _ := r.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))

	if text == "" {
		return defaultYes
	}
	return text == "y" || text == "yes"
}
