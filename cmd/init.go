package cmd

import (
	"os"
	"path"

	"github.com/smtdfc/contractor/helpers"
	"github.com/spf13/cobra"
)

type InitData struct{}

func InitCommandFn(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	paths := map[string]string{
		"templates/project/contractor.config.json.tmpl":  path.Join(cwd, "./contractor.config.json"),
		"templates/project/contract/model.contract.tmpl": path.Join(cwd, "./contract/model.contract"),
	}

	for t, v := range paths {
		err := helpers.RenderTemplateFromFile(t, v, &InitData{})
		if err != nil {
			return err
		}
	}

	return nil
}

var InitCommand = &cobra.Command{
	Use:   "init",
	Short: "Init project",
	Long:  "Init project",
	RunE:  InitCommandFn,
}
