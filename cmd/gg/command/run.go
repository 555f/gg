package command

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/gg"

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
	Short: "Starts code generation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		packageNames := viper.GetStringSlice("packages")
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
		cmd.Printf(yellow("Packages: %s"), strings.Join(escaped, ","))

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
			Dir:        wd,
			Env:        os.Environ(),
			BuildFlags: []string{"-tags=gg"},
		}

		pkgs, err := stdpackages.Load(cfg, escaped...)
		if err != nil {
			return
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

		// if len(files) > 0 {
		// 	cmd.Printf(green("\n\nfiles was generated:\n"))
		// 	for _, f := range files {
		// 		data, err := f.Bytes()
		// 		if err != nil {
		// 			cmd.Printf("%s %s %s: %s", red("𐄂"), red("error during file generation"), yellow(f.Filepath()), red(err))
		// 			continue
		// 		}
		// 		dirPath := filepath.Dir(f.Filepath())
		// 		if err := os.MkdirAll(dirPath, 0700); err != nil {
		// 			cmd.Printf("%s %s %s: %s", red("𐄂"), red("error when creating a directory"), yellow(dirPath), red(err))
		// 			continue
		// 		}
		// 		if err := os.WriteFile(f.Filepath(), data, 0700); err != nil {
		// 			cmd.Printf("%s %s %s: %s", red("𐄂"), red("error when creating a file"), yellow(f.Filepath()), red(err))
		// 			continue
		// 		}
		// 		cmd.Println(green("✓"), f.Filepath())
		// 	}
		// } else {
		// 	cmd.Printf("\nnothing generation\n")
		// }
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceP("packages", "p", nil, "Scan packages")
	_ = viper.BindPFlag("packages", runCmd.Flags().Lookup("packages"))
}
