package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

var targetDir string

var rootCmd = &cobra.Command{
	Use:   "npmdepcopy",
	Short: "Copy NodeJS module dependencies for production to a target directory.",
	RunE:  depCopy,
}

func init() {
	rootCmd.Flags().StringVarP(&targetDir, "out", "o", "", "Target output directory to copy modules to (required)")
	rootCmd.MarkFlagRequired("out")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func depCopy(cmd *cobra.Command, args []string) error {

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory\n%f", err)
	}

	out, err := exec.Command("npm", "ls", "--prod", "--parseable").Output()
	if err != nil {
		return fmt.Errorf("Failed to get list of dependencies from npm, is it installed?\n%f", err)
	}

	packages := strings.Split(string(out), "\n")
	for _, packDir := range packages {
		packDir = strings.Replace(packDir, dir, "", -1)
		if packDir == "" {
			continue
		}

		packName := strings.Replace(packDir, "/node_modules/", "", -1)

		targetDir = strings.TrimRight(targetDir, "/")

		fmt.Printf("Copying module '%s' to '%s' ... ", packName, targetDir)
		if err := copy.Copy(dir+packDir, targetDir+"/"+packName); err != nil {
			fmt.Printf("%s\n\n", aurora.Red("FAILED"))
			return err
		}

		fmt.Printf("%s\n", aurora.Green("SUCCEEDED"))
	}

	fmt.Println(aurora.Green("\nOperation completed successfully!"))

	return nil
}

func getCwd() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func parsePackageList(packList string) ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return []string{}, fmt.Errorf("failed to get current working directory: %f", err)
	}

	packages := strings.Split(packList, "\n")
	filtered := []string{}
	for _, pack := range packages {
		pack = strings.Replace(pack, dir, "", -1)
		if pack != "" {
			filtered = append(filtered, pack)
		}
	}

	return filtered, nil
}
