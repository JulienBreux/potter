package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	cmdRoot "github.com/julienbreux/potter/internal/potter/command/root"
	cmdVersion "github.com/julienbreux/potter/internal/potter/command/version"
	ver "github.com/julienbreux/potter/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgPathFile string
	cfgSubPath  = ".config/potter/"
	cfgFile     = "potter.yml"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgPathFile != "" {
		viper.SetConfigFile(cfgPathFile)
	} else {
		viper.AddConfigPath(cfgPath())
		viper.SetConfigType("yml")
		viper.SetConfigName(cfgFile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// Execute executes command
func New(in io.Reader, out, errIO io.Writer, args ...string) *cobra.Command {
	defer func() {
		if r := recover(); r != nil {
			// TODO: Improve error message color
			fmt.Println("Internal Potter error")
			// TODO: Add logger at debug level
			// TODO: Add "tips" option
			// TODO: Get URL from outside
			fmt.Println("➡ Please report here: https://github.com/julienbreux/potter/issues/new?labels=bug")
			os.Exit(1)
		}
	}()

	// Cobra initialization
	cobra.OnInitialize(initConfig)

	// Create root command
	cmd := &cobra.Command{
		Use:     "potter",
		Short:   "Potter is a magical artifact to enchant the world of containers",
		Version: ver.Version,
		Run:     cmdRoot.Run,
	}

	cmd.SetIn(in)
	cmd.SetOut(out)
	cmd.SetErr(errIO)
	cmd.SetArgs(args)

	// Add flags
	flags(cmd)

	// Add subcommands
	cmd.AddCommand(cmdVersion.New())

	return cmd
}

func flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&cfgPathFile, "config", "c", filepath.Join(cfgPath(), cfgFile), "configuration file")
	cmd.PersistentFlags().StringP("output", "o", "", "one of '', 'yaml' or 'json'.")
}

func cfgPath() string {
	// Home directory
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return filepath.Join(homeDir, cfgSubPath)
}
