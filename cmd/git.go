package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/bono/loadstar/internal/git"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Manage git remote integration settings",
	Long: `Configure the git remote used by 'loadstar checkpoint' for automatic push.
Credentials are stored in .loadstar/COMMON/git_config.json (gitignored).

Subcommands:
  set     Configure remote URL, branch, and PAT
  status  Show current remote configuration and repo state
  unset   Remove remote configuration

Examples:
  loadstar git set https://github.com/user/repo.git --user aeolusk --email u@example.com --token <PAT>
  loadstar git status
  loadstar git unset`,
}

var gitSetCmd = &cobra.Command{
	Use:   "set [URL]",
	Short: "Configure git remote for loadstar checkpoint push",
	Long: `Save git remote settings so that 'loadstar checkpoint' can push automatically.
The PAT is stored in plaintext in .loadstar/COMMON/git_config.json (gitignored).

Examples:
  loadstar git set https://github.com/user/repo.git
  loadstar git set https://github.com/user/repo.git --branch main --user aeolusk --email u@example.com --token ghp_xxx`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remoteURL := args[0]
		branch, _ := cmd.Flags().GetString("branch")
		user, _ := cmd.Flags().GetString("user")
		email, _ := cmd.Flags().GetString("email")
		token, _ := cmd.Flags().GetString("token")

		if branch == "" {
			branch = "main"
		}

		gc := git.NewClient(fs.Root)
		if err := gc.SetRemote(remoteURL, branch, user, email, token); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("remote configured: %s (branch: %s)\n", remoteURL, branch)
		if token != "" {
			fmt.Println("PAT saved to .loadstar/COMMON/git_config.json")
		}
	},
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current git integration status",
	Run: func(cmd *cobra.Command, args []string) {
		gc := git.NewClient(fs.Root)
		info, err := gc.GetStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if info.RemoteURL == "" {
			fmt.Fprintln(w, "Remote:\t(not configured — run `loadstar git set <URL>`)")
		} else {
			fmt.Fprintf(w, "Remote:\t%s\n", info.RemoteURL)
		}
		if info.Branch == "" {
			fmt.Fprintln(w, "Branch:\t(unknown)")
		} else {
			fmt.Fprintf(w, "Branch:\t%s\n", info.Branch)
		}
		if info.LatestHash == "" {
			fmt.Fprintln(w, "Latest commit:\t(none)")
		} else {
			fmt.Fprintf(w, "Latest commit:\t%s\n", info.LatestHash[:8])
		}
		fmt.Fprintf(w, "Uncommitted files:\t%d\n", info.UncommittedFiles)
		w.Flush()
	},
}

var gitUnsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Remove git remote configuration",
	Run: func(cmd *cobra.Command, args []string) {
		gc := git.NewClient(fs.Root)
		if err := gc.UnsetRemote(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("git remote configuration removed")
	},
}

func init() {
	gitCmd.AddCommand(gitSetCmd)
	gitCmd.AddCommand(gitStatusCmd)
	gitCmd.AddCommand(gitUnsetCmd)

	gitSetCmd.Flags().String("branch", "main", "Branch to push to (default: main)")
	gitSetCmd.Flags().String("user", "", "Git author name")
	gitSetCmd.Flags().String("email", "", "Git author email")
	gitSetCmd.Flags().String("token", "", "Personal Access Token for push authentication")
}
