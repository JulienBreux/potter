package root

import (
	"github.com/julienbreux/potter/internal/potter/webui"
	"github.com/julienbreux/potter/pkg/version"
	"github.com/spf13/cobra"
)

// run returns the command
func Run(cmd *cobra.Command, args []string) {
	_ = webui.New(version.Version).Run()
}
