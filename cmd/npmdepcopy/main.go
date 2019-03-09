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

var inputDir string
var outputDir string

var rootCmd = &cobra.Command{
	Use:   "npmdepcopy",
	Short: "Copy NodeJS module dependencies for production to a target directory.",
	Run: func(cmd *cobra.Command, args []string) {
		err := depCopy(cmd, args)
		if err != nil {
			fmt.Println(aurora.Red(err.Error()))
		}
	},
}

func init() {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	rootCmd.Flags().StringVarP(&inputDir, "in", "i", ".", "The input directory where the package.json and node_modules live")
	rootCmd.Flags().StringVarP(&outputDir, "out", "o", "", "Target output directory to copy modules to (required)")
	rootCmd.MarkFlagRequired("out")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func depCopy(cmd *cobra.Command, args []string) (err error) {

	input, err := filepath.Abs(inputDir)
	if err != nil {
		return fmt.Errorf("unable to generate input path from option or current working directory\n%f", err)
	}

	stat, err := os.Stat(input)
	if err != nil || !stat.Mode().IsDir() {
		return fmt.Errorf("given input path does not point to a valid directory\n%f", err)
	}

	npmi := exec.Command("npm", "i")
	npmi.Dir = input
	_, err = npmi.Output()
	if err != nil {
		return fmt.Errorf("failed to install dependencies for copying\n%f", err)
	}

	npmls := exec.Command("npm", "ls", "--prod", "--parseable")
	npmls.Dir = input
	result, err := npmls.Output()
	if err != nil {
		return fmt.Errorf("failed to get list of dependencies from npm, is it installed?\n%f", err)
	}

	out, err := filepath.Abs(strings.TrimRight(outputDir, "/"))
	if err != nil {
		return fmt.Errorf("failed to create absolute filepath for output directory\n%f", err)
	}

	packages := strings.Split(string(result), "\n")
	for _, packDir := range packages {
		packDir = strings.Replace(packDir, input, "", -1)
		if packDir == "" {
			continue
		}

		packName := strings.Replace(packDir, "/node_modules/", "", -1)
		fmt.Println(packName)

		fmt.Printf("Copying module '%s' to '%s' ... ", packName, out)
		pinput := filepath.Join(input, packDir)
		poutput := filepath.Join(out, packName)

		if err := copy.Copy(pinput, poutput); err != nil {
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
