package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "Show the version, commit, and build date of the QuickSpin CLI",
		Run:   runVersion,
	}

	return cmd
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("qspin version %s\n", Version)
	fmt.Printf("commit: %s\n", Commit)
	fmt.Printf("built: %s\n", Date)
}
