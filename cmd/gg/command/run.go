package command

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
	stdpackages "golang.org/x/tools/go/packages"
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
		noSelfUpdate := viper.GetBool("no-selfupdate")

		if !noSelfUpdate {
			var updater = &selfupdate.Updater{
				CurrentVersion: cmd.Root().Version,          // Manually update the const, or set it using `go build -ldflags="-X main.VERSION=<newver>" -o hello-updater src/hello-updater/main.go`
				ApiURL:         "http://51.250.88.10:8081/", // The server hosting `$CmdName/$GOOS-$ARCH.json` which contains the checksum for the binary
				BinURL:         "http://51.250.88.10:8081/", // The server hosting the zip file containing the binary application which is a fallback for the patch method
				Dir:            "update/",                   // The directory created by the app when run which stores the cktime file
				CmdName:        "",                          // The app name which is appended to the ApiURL to look for an update
				ForceCheck:     true,                        // For this example, always check for an update unless the version is "dev"
			}
			version, err := updater.UpdateAvailable()
			if err != nil {
				cmd.Printf(red("Check update failed %v\n"), err)
				return
			}
			if version != "" {
				cmd.Print(green("Update available!!!\n"))

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
		isDebug := viper.GetBool("debug")
		configFile := viper.ConfigFileUsed()
		configFile = filepath.FromSlash(configFile)

		wd := filepath.Dir(configFile)
		wdAbs, _ := filepath.Abs(wd)

		plugins := viper.GetStringMap("plugins")

		escaped := make([]string, len(packageNames))
		for i := range packageNames {
			escaped[i] = "pattern=" + packageNames[i]
		}

		cmd.Printf(yellow("Version: %s\n"), cmd.Root().Version)
		cmd.Printf(yellow("Workdir: %s\n"), wdAbs)
		cmd.Printf(yellow("Config file: %s\n"), configFile)
		cmd.Printf(yellow("Packages: %s\n"), strings.Join(escaped, ","))

		cfg := &stdpackages.Config{
			ParseFile: func(fSet *token.FileSet, filename string, src []byte) (*ast.File, error) {
				return parser.ParseFile(fSet, filename, src, parser.AllErrors|parser.ParseComments)
			},
			Mode: stdpackages.NeedDeps |
				stdpackages.NeedSyntax |
				stdpackages.NeedTypesInfo |
				stdpackages.NeedTypes |
				stdpackages.NeedTypesSizes |
				stdpackages.NeedImports |
				stdpackages.NeedName |
				stdpackages.NeedModule |
				stdpackages.NeedFiles |
				stdpackages.NeedCompiledGoFiles,
			Dir:        wdAbs,
			Env:        os.Environ(),
			BuildFlags: []string{"-tags=gg"},
		}

		pkgs, err := stdpackages.Load(cfg, escaped...)
		if err != nil {
			cmd.Printf("\n\n%s\n", red(err))
			os.Exit(1)
		}

		if isDebug {
			cmd.Printf(yellow("Found packages: %d\n"), len(pkgs))
		}

		var foundPkgErr bool
		for _, pkg := range pkgs {
			if len(pkg.Errors) > 0 {
				foundPkgErr = true
				for _, err := range pkg.Errors {
					cmd.Printf("\n\n%s\n", red(err))
				}
			}
		}
		if foundPkgErr {
			os.Exit(1)
		}

		var isExitApp bool

		result, err := gg.Run(cmd.Root().Version, wdAbs, pkgs, plugins)
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
				cmd.Println(green("✓"), r.Filepath)
			}
		} else {
			cmd.Printf("\nnothing generation\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceP("packages", "p", nil, "Scan packages")
	_ = viper.BindPFlag("packages", runCmd.Flags().Lookup("packages"))

	runCmd.Flags().BoolP("debug", "d", false, "Debug mode")
	_ = viper.BindPFlag("debug", runCmd.Flags().Lookup("debug"))

	runCmd.Flags().BoolP("no-selfupdate", "s", false, "Disable self-update")
	_ = viper.BindPFlag("no-selfupdate", runCmd.Flags().Lookup("no-selfupdate"))
}
