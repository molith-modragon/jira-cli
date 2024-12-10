package versions

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/versions/list"
)

const helpText = `Project manages Jira projects. See available commands below.`

// NewCmdProject is a project command.
func NewCmdVersions() *cobra.Command {
	cmd := cobra.Command{
		Use:         "versions",
		Short:       "versions manages project versions",
		Long:        helpText,
		Aliases:     []string{"releases"},
		Annotations: map[string]string{"cmd:main": "true"},
		RunE:        versions,
	}

	cmd.AddCommand(list.NewCmdList())

	return &cmd
}

func versions(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
