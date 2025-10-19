package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run:   printVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion(cmd *cobra.Command, args []string) {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string `json:"release"`
		BuildDate string `json:"build_date"`
		GitHash   string `json:"git_hash"`
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
		os.Exit(1)
	}
}
