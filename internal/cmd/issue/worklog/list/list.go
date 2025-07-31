package list

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/query"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `List displays all worklogs for an issue with their complete attributes.

When --tempo flag is used, the command will also fetch custom Tempo worklog attributes 
(requires Tempo plugin and API configuration).`
	examples = `$ jira issue worklog list

# List worklogs for a specific issue
$ jira issue worklog list ISSUE-1

# List worklogs in plain output format
$ jira issue worklog list ISSUE-1 --plain

# List worklogs with Tempo custom attributes
$ jira issue worklog list ISSUE-1 --tempo

# List worklogs with Tempo attributes in plain format
$ jira issue worklog list ISSUE-1 --tempo --plain`
)

// NewCmdWorklogList is a worklog list command.
func NewCmdWorklogList() *cobra.Command {
	cmd := cobra.Command{
		Use:     "list ISSUE-KEY",
		Short:   "List all worklogs for an issue",
		Long:    helpText,
		Example: examples,
		Annotations: map[string]string{
			"help:args": "ISSUE-KEY\tIssue key to list worklogs for, eg: ISSUE-1",
		},
		Run: list,
	}

	cmd.Flags().Bool("plain", false, "Display output in plain text")
	cmd.Flags().Bool("tempo", false, "Include Tempo custom attributes (requires Tempo API configuration)")

	return &cmd
}

func list(cmd *cobra.Command, args []string) {
	params := parseArgsAndFlags(args, cmd.Flags())
	client := api.DefaultClient(params.debug)

	lc := listCmd{
		client: client,
		params: params,
	}

	cmdutil.ExitIfError(lc.setIssueKey())
	cmdutil.ExitIfError(lc.run())
}

type listCmd struct {
	client *jira.Client
	params *listParams
}

type listParams struct {
	issueKey string
	debug    bool
	plain    bool
	tempo    bool
}

func parseArgsAndFlags(args []string, flags query.FlagParser) *listParams {
	var issueKey string

	nArgs := len(args)
	if nArgs >= 1 {
		issueKey = args[0]
	}

	debug, err := flags.GetBool("debug")
	cmdutil.ExitIfError(err)

	plain, err := flags.GetBool("plain")
	cmdutil.ExitIfError(err)

	tempo, err := flags.GetBool("tempo")
	cmdutil.ExitIfError(err)

	return &listParams{
		issueKey: issueKey,
		debug:    debug,
		plain:    plain,
		tempo:    tempo,
	}
}

func (lc *listCmd) setIssueKey() error {
	if lc.params.issueKey != "" {
		return nil
	}

	var ans string
	qs := &survey.Question{
		Name:     "issueKey",
		Prompt:   &survey.Input{Message: "Issue key"},
		Validate: survey.Required,
	}
	if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
		return err
	}
	lc.params.issueKey = cmdutil.GetJiraIssueKey(viper.GetString("project.key"), ans)

	return nil
}

func (lc *listCmd) run() error {
	s := cmdutil.Info(fmt.Sprintf("Fetching worklogs for issue %s...", lc.params.issueKey))
	defer s.Stop()

	if lc.params.tempo {
		// Use enhanced API with Tempo support
		worklogsWithTempo, err := lc.client.GetIssueWorklogsWithTempo(lc.params.issueKey, true)
		if err != nil {
			return err
		}

		s.Stop()

		if len(worklogsWithTempo) == 0 {
			cmdutil.Failed("No worklogs found for issue %s", lc.params.issueKey)
			return nil
		}

		view.PrintWorklogsWithTempo(worklogsWithTempo, lc.params.plain)
	} else {
		// Use standard API
		worklogList, err := lc.client.GetIssueWorklogs(lc.params.issueKey)
		if err != nil {
			return err
		}

		s.Stop()

		if len(worklogList.Worklogs) == 0 {
			cmdutil.Failed("No worklogs found for issue %s", lc.params.issueKey)
			return nil
		}

		view.PrintWorklogs(worklogList.Worklogs, lc.params.plain)
	}

	return nil
}