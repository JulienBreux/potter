package root

import (
	"github.com/JulienBreux/potter/internal/potter/webui"
	"github.com/spf13/cobra"
)

// run returns the command
func Run(cmd *cobra.Command, args []string) {
	_ = webui.New().Run()
}
