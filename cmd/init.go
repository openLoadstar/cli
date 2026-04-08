package cmd

import (
	"fmt"
	"os"

	"github.com/bono/loadstar/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise .loadstar/ structure in the current directory",
	Long: `Initialise the .loadstar/ directory structure.

Creates MAP/, WAYPOINT/, COMMON/, and .clionly/ subdirectories,
along with an initial M://root Map file.

Git initialisation should be done separately with 'git init'.

Examples:
  loadstar init`,
	// Override PersistentPreRun — init must work before .loadstar/ exists.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: cannot determine working directory:", err)
			os.Exit(1)
		}

		if err := storage.Init(cwd); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to initialise .loadstar/:", err)
			os.Exit(1)
		}
		fmt.Println("initialised .loadstar/ structure")
		fmt.Println("done. use 'git init' separately if needed.")
	},
}
