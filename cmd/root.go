package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	basesrepo "github.com/kawana77b/tsconfig-template/internal/bases_repo"
	"github.com/kawana77b/tsconfig-template/internal/prompt"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

var Version string = "0.0.0"

var TemplateFs embed.FS

type rootProps struct {
	useDefault bool
	saveDir    string
}

var rootP = &rootProps{}

var rootCmd = &cobra.Command{
	Use:   "tsconfig-template",
	Args:  cobra.NoArgs,
	Short: "Select tsconfig.json and output its contents.",
	Long: `Select tsconfig.json and output its contents.

This program's templates uses the following resources:

- tsconfig/bases https://github.com/tsconfig/bases
  Copyright (c) Microsoft Corporation.
	`,
	PreRunE: preRunRoot,
	RunE:    runRoot,
}

func Execute() {
	rootCmd.Version = Version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&rootP.useDefault, "default", "d", false, "use default (recommended) template")
}

func preRunRoot(cmd *cobra.Command, args []string) error {
	// The destination of tsconfig is set to the current current folder.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	rootP.saveDir = wd

	return nil
}

func runRoot(cmd *cobra.Command, args []string) error {
	repo := basesrepo.NewBasesRepo(TemplateFs)
	files, err := repo.TemplateFiles()
	if err != nil {
		return err
	}

	fileMap, names := basesrepo.GetFileMap(files)

	var choice string
	if rootP.useDefault {
		choice = basesrepo.DEFAULT_JSON_FILENAME
	}

	if choice == "" {
		ans, err := prompt.NewSelectQuestion().Ask("Select a template", names)
		if err != nil {
			// If the user cancels, return nil because you want to finish without displaying any error.
			return nil
		}
		choice = ans
	}
	path, ok := fileMap[choice]
	if !ok {
		return fmt.Errorf("file not found: %s", choice)
	}

	bytes, err := repo.ReadFile(path)
	if err != nil {
		return err
	}

	savePath := filepath.Join(rootP.saveDir, "tsconfig.json")
	os.WriteFile(savePath, bytes, 0644)

	fmt.Printf("%s %s\n", aurora.Cyan("created:"), savePath)

	return nil
}
