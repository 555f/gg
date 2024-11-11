package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/selfupdate"
	"github.com/manifoldco/promptui"

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

const (
	defaultSelfUpdateAPI    = "https://gg.lobchuk.ru"
	defaultSelfUpdateBinary = "https://gg.lobchuk.ru"
)

// runCmd represents the init command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start code generation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		noSelfUpdate := viper.GetBool("no-selfupdate")
		if cmd.Root().Version != "dev" && !noSelfUpdate {
			var updater = &selfupdate.Updater{
				CurrentVersion: cmd.Root().Version,                // Manually update the const, or set it using `go build -ldflags="-X main.VERSION=<newver>" -o hello-updater src/hello-updater/main.go`
				ApiURL:         viper.GetString("selfupdate-api"), // The server hosting `$CmdName/$GOOS-$ARCH.json` which contains the checksum for the binary
				BinURL:         viper.GetString("selfupdate-bin"), // The server hosting the zip file containing the binary application which is a fallback for the patch method
				Dir:            "update/",                         // The directory created by the app when run which stores the cktime file
				CmdName:        "",                                // The app name which is appended to the ApiURL to look for an update
				ForceCheck:     true,                              // For this example, always check for an update unless the version is "dev"
			}
			version, err := updater.UpdateAvailable()
			if err != nil {
				cmd.Printf(red("Check update failed %v\n"), err)
				return
			}
			if version != "" {
				cmd.Print(green("Update version ") + yellow(version) + green(" available!!!\n"))

				prompt := promptui.Prompt{
					Label:     "Do you have update",
					IsConfirm: true,
				}
				result, _ := prompt.Run()
				if result == "y" || result == "yes" {
					err := updater.Run()
					if err != nil {
						log.Println("Failed to update app:", err)
					}
					cmd.Printf(green("Update to latest version: %s success, run command again.\n"), updater.Info.Version)
					return
				}
			}
		}

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

	runCmd.Flags().BoolP("no-selfupdate", "s", false, "Self update disable")
	_ = viper.BindPFlag("no-selfupdate", runCmd.Flags().Lookup("no-selfupdate"))

	runCmd.Flags().StringP("selfupdate-api", "a", defaultSelfUpdateAPI, "The server hosting `$CmdName/$GOOS-$ARCH.json` which contains the checksum for the binary")
	_ = viper.BindPFlag("selfupdate-api", runCmd.Flags().Lookup("selfupdate-api"))

	runCmd.Flags().StringP("selfupdate-bin", "b", defaultSelfUpdateBinary, "The server hosting the zip file containing the binary application which is a fallback for the patch method")
	_ = viper.BindPFlag("selfupdate-bin", runCmd.Flags().Lookup("selfupdate-bin"))

	runCmd.Flags().StringToStringP("plugins", "l", map[string]string{}, "")
	_ = viper.BindPFlag("plugins", runCmd.Flags().Lookup("plugins"))
}
