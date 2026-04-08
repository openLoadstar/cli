package cmd

import (
	"fmt"
	"os"

	"github.com/bono/loadstar/internal/core"
	"github.com/bono/loadstar/internal/storage"
	"github.com/spf13/cobra"
)

// svc and fs are initialised once at PersistentPreRun and reused by all subcommands.
var (
	svc *core.ElementService
	fs  *storage.FS
)

var rootCmd = &cobra.Command{
	Use:   "loadstar",
	Short: "LOADSTAR - Project metadata and waypoint management CLI",
	Long: `LOADSTAR CLI for managing project metadata and waypoints.

Working directory:
  loadstar searches for .loadstar/ starting from the current directory and
  walking up the directory tree. Run it from anywhere inside your project.

  If .loadstar/ is not found, a new one is auto-initialised in the current
  directory — so run from inside your project root to avoid accidental init.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: cannot determine working directory:", err)
			os.Exit(1)
		}

		// Auto-initialise .loadstar/ if not found anywhere up the tree.
		root, err := storage.FindRoot(cwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to initialise .loadstar/:", err)
			os.Exit(1)
		}

		fs = storage.NewFS(root)
		svc = core.NewElementService(fs)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(todoCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(validateCmd)
}
