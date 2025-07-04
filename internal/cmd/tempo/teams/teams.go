package teams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/tempo"
	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

const helpText = `Teams command retrieves teams from Tempo.`

// NewCmdTeams is a teams command.
func NewCmdTeams() *cobra.Command {
	cmd := cobra.Command{
		Use:     "teams",
		Short:   "Retrieve teams from Tempo",
		Long:    helpText,
		Aliases: []string{"team"},
		RunE:    teams,
	}

	cmd.Flags().Bool("plain", false, "Display output in plain text format")
	cmd.Flags().String("format", "table", "Output format (table, json)")

	return &cmd
}

func teams(cmd *cobra.Command, _ []string) error {
	debug := viper.GetBool("debug")
	plain, err := cmd.Flags().GetBool("plain")
	if err != nil {
		return err
	}
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return err
	}

	client := api.DefaultClient(debug)
	tempoClient := tempo.NewClient(client)

	teams, err := tempoClient.GetTeams()
	if err != nil {
		return err
	}

	if format == "json" {
		return displayAsJSON(teams)
	}

	if plain {
		return displayPlainText(teams)
	}

	return displayTable(teams)
}

func displayAsJSON(teams []*tempo.Team) error {
	data, err := json.MarshalIndent(teams, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func displayPlainText(teams []*tempo.Team) error {
	for _, team := range teams {
		fmt.Printf("ID: %d, Name: %s, Summary: %s\n", team.ID, team.Name, team.Summary)
	}
	return nil
}

func displayTable(teams []*tempo.Team) error {
	if len(teams) == 0 {
		cmdutil.Failed("No teams found")
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 8, 1, '\t', 0)

	// Print header
	fmt.Fprintf(w, "ID\tNAME\tSUMMARY\n")

	// Print teams
	for _, team := range teams {
		fmt.Fprintf(w, "%d\t%s\t%s\n", team.ID, team.Name, team.Summary)
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return tui.PagerOut(buf.String())
}