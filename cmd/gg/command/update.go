package command

import (
	"log"

	"github.com/555f/selfupdate"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultSelfUpdateAPI    = "https://gg.lobchuk.ru"
	defaultSelfUpdateBinary = "https://gg.lobchuk.ru"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Self update",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Root().Version != "dev" {
			var updater = &selfupdate.Updater{
				CurrentVersion: cmd.Root().Version,     // Manually update the const, or set it using `go build -ldflags="-X main.VERSION=<newver>" -o hello-updater src/hello-updater/main.go`
				ApiURL:         viper.GetString("api"), // The server hosting `$CmdName/$GOOS-$ARCH.json` which contains the checksum for the binary
				BinURL:         viper.GetString("bin"), // The server hosting the zip file containing the binary application which is a fallback for the patch method
				Dir:            "update/",              // The directory created by the app when run which stores the cktime file
				CmdName:        "",                     // The app name which is appended to the ApiURL to look for an update
				ForceCheck:     true,                   // For this example, always check for an update unless the version is "dev"
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
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("api", "a", defaultSelfUpdateAPI, "The server hosting `$CmdName/$GOOS-$ARCH.json` which contains the checksum for the binary")
	_ = viper.BindPFlag("aåpi", updateCmd.Flags().Lookup("api"))

	updateCmd.Flags().StringP("bin", "b", defaultSelfUpdateBinary, "The server hosting the zip file containing the binary application which is a fallback for the patch method")
	_ = viper.BindPFlag("bin", updateCmd.Flags().Lookup("bin"))
}
