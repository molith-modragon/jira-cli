package list

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `List lists releases in a given project.`
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List lists releases in a project",
		Long:    helpText,
		Aliases: []string{"lists", "ls"},
		Run:     List,
	}
}

// List displays a list view.
func List(cmd *cobra.Command, _ []string) {
	loadList(cmd)
}

func loadList(cmd *cobra.Command) {
	project := viper.GetString("project.key")

	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	versions, total, err := func() ([]*jira.Version, int, error) {
		s := cmdutil.Info("Fetching versions...")
		defer s.Stop()

		versions, err := api.DefaultClient(debug).Version(project)
		if err != nil {
			return nil, 0, err
		}
		return versions, len(versions), nil
	}()
	cmdutil.ExitIfError(err)

	if total == 0 {
		fmt.Println()
		cmdutil.Failed("No result found for given query in project %q", project)
		return
	}

	v := view.NewVersion(versions)

	cmdutil.ExitIfError(v.Render())
}
