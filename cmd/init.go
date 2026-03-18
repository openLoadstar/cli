package cmd

import (
	"fmt"
	"os"

	"github.com/bono/loadstar/internal/git"
	"github.com/bono/loadstar/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise .loadstar/ structure and git repository in the current directory",
	Long: `Initialise the .loadstar/ directory structure and a local git repository.

Optionally configure a remote repository so that 'loadstar checkpoint' can push
automatically. The PAT is stored in .loadstar/COMMON/git_config.json (plaintext
for now; encryption is planned for a future release).

Examples:
  loadstar init
  loadstar init --remote https://github.com/aeolusk/repo.git --branch main --user aeolusk --email user@example.com --token <PAT>`,
	// Override PersistentPreRun — init must work before .loadstar/ exists.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: cannot determine working directory:", err)
			os.Exit(1)
		}

		// 1. Initialise .loadstar/ directory structure
		if err := storage.Init(cwd); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to initialise .loadstar/:", err)
			os.Exit(1)
		}
		fmt.Println("initialised .loadstar/ structure")

		// 2. Initialise git repository
		gc := git.NewClient(cwd)
		if err := gc.Init(); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to initialise git repository:", err)
			os.Exit(1)
		}
		fmt.Println("initialised git repository")

		// 3. Configure remote if --remote flag provided
		remote, _ := cmd.Flags().GetString("remote")
		if remote != "" {
			branch, _ := cmd.Flags().GetString("branch")
			user, _ := cmd.Flags().GetString("user")
			email, _ := cmd.Flags().GetString("email")
			token, _ := cmd.Flags().GetString("token")

			if branch == "" {
				branch = "main"
			}

			if err := gc.SetRemote(remote, branch, user, email, token); err != nil {
				fmt.Fprintln(os.Stderr, "error: failed to configure remote:", err)
				os.Exit(1)
			}
			fmt.Printf("remote configured: %s (branch: %s)\n", remote, branch)
			if token != "" {
				fmt.Println("PAT saved to .loadstar/COMMON/git_config.json")
			}
		}

		fmt.Println("done. run `loadstar checkpoint -m \"initial\"` to create the first commit.")
	},
}

func init() {
	initCmd.Flags().String("remote", "", "Remote repository URL (e.g. https://github.com/user/repo.git)")
	initCmd.Flags().String("branch", "main", "Branch to push to (default: main)")
	initCmd.Flags().String("user", "", "Git author name")
	initCmd.Flags().String("email", "", "Git author email")
}
