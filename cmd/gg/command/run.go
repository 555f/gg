package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/gg"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

// runCmd represents the init command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start code generation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		packageNames := viper.GetStringSlice("packages")
		configFile := viper.ConfigFileUsed()
		configFile = filepath.FromSlash(configFile)

		var (
			wdAbs string
		)

		if configFile != "" {
			wdAbs, _ = filepath.Abs(filepath.Dir(configFile))
		} else {
			wdAbs, _ = filepath.Abs(viper.GetString("wd"))
		}

		for i, pkgPath := range packageNames {
			packageNames[i] = filepath.Join(wdAbs, pkgPath)
		}

		pluginOpts := viper.GetStringMapString("plugins")

		cmd.Printf(yellow("Version: %s\n"), cmd.Root().Version)
		cmd.Printf(yellow("Workdir: %s\n"), wdAbs)
		if configFile != "" {
			cmd.Printf(yellow("Config file: %s\n"), configFile)
		}

		var isExitApp bool

		result, err := gg.Run(cmd.Root().Version, wdAbs, packageNames, pluginOpts, true)
		if err != nil {
			if merr, ok := err.(*multierror.Error); ok {
				merr.ErrorFormat = func(es []error) string {
					errorPoints := make([]string, 0, len(es))
					warningPoints := make([]string, 0, len(es))
					for _, err := range es {
						if errors.IsWarning(err) {
							warningPoints = append(warningPoints, fmt.Sprintf("* %s", yellow(err)))
						} else {
							isExitApp = true
							errorPoints = append(errorPoints, fmt.Sprintf("* %s", red(err)))
						}
					}
					var text string
					if len(errorPoints) > 0 {
						text += fmt.Sprintf(
							"\n\n%d error(s) occurred:\n\t%s\n\n",
							len(errorPoints), strings.Join(errorPoints, "\n\t"))
					}
					if len(warningPoints) > 0 {
						text += fmt.Sprintf(
							"\n\n%d warning(s) occurred:\n\t%s\n\n",
							len(warningPoints), strings.Join(warningPoints, "\n\t"))
					}
					return text
				}
			}
			cmd.Println(err)
		}
		if isExitApp {
			os.Exit(1)
		}
		if len(result) > 0 {
			cmd.Printf(green("\n\nfiles was generated:\n"))
			for _, r := range result {
				cmd.Println(green("✓"), r.File.Path())
			}
		} else {
			cmd.Printf("\nnothing generation\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("wd", "w", "", "Work dir")
	_ = viper.BindPFlag("wd", runCmd.Flags().Lookup("wd"))

	runCmd.Flags().StringSliceP("packages", "p", nil, "Scan packages")
	_ = viper.BindPFlag("packages", runCmd.Flags().Lookup("packages"))

	runCmd.Flags().BoolP("debug", "d", false, "Debug mode")
	_ = viper.BindPFlag("debug", runCmd.Flags().Lookup("debug"))

	runCmd.Flags().StringToStringP("plugins", "l", map[string]string{}, "")
	_ = viper.BindPFlag("plugins", runCmd.Flags().Lookup("plugins"))
}
