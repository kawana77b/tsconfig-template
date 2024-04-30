package cmd

import (
	"fmt"

	basesrepo "github.com/kawana77b/tsconfig-template/internal/bases_repo"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List template types",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := basesrepo.NewBasesRepo(TemplateFs)
		files, err := repo.TemplateFiles()
		if err != nil {
			return err
		}

		for _, f := range files {
			fmt.Println(f.Name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
