package tempo

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/tempo/teams"
)

const helpText = `Tempo manages Tempo resources. See available commands below.`

// NewCmdTempo is a tempo command.
func NewCmdTempo() *cobra.Command {
	cmd := cobra.Command{
		Use:         "tempo",
		Short:       "Tempo manages Tempo resources",
		Long:        helpText,
		Annotations: map[string]string{"cmd:main": "true"},
		RunE:        tempo,
	}

	cmd.AddCommand(teams.NewCmdTeams())

	return &cmd
}

func tempo(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}