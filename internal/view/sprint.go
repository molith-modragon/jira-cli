package view

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

// SprintIssueFunc provides issues in the sprint.
type SprintIssueFunc func(boardID, sprintID int) []*jira.Issue

// SprintList is a list view for sprints.
type SprintList struct {
	Project string
	Board   string
	Server  string
	Data    []*jira.Sprint
	Issues  SprintIssueFunc
	Display DisplayFormat

	issueCache map[string]tui.TableData
}

// Render renders the sprint explorer view.
func (sl SprintList) Render() error {
	data := sl.data()

	view := tui.NewPreview(
		tui.WithPreviewFooterText(
			fmt.Sprintf(
				"Showing %d results from board \"%s\" of project \"%s\"",
				len(sl.Data), sl.Board, sl.Project,
			),
		),
		tui.WithInitialText(helpText),
		tui.WithContentTableOpts(tui.WithSelectedFunc(navigate(sl.Server))),
	)

	return view.Render(data)
}

// RenderInTable renders the list in table view.
func (sl SprintList) RenderInTable() error {
	if sl.Display.Plain {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
		return sl.renderPlain(w)
	}

	data := sl.tableData()

	view := tui.NewTable(
		tui.WithColPadding(colPadding),
		tui.WithMaxColWidth(maxColWidth),
		tui.WithTableFooterText(
			fmt.Sprintf(
				"Showing %d results from board \"%s\" of project \"%s\"",
				len(sl.Data), sl.Board, sl.Project,
			),
		),
	)

	return view.Render(data)
}

// renderPlain renders the issue in plain view.
func (sl SprintList) renderPlain(w io.Writer) error {
	return renderPlain(w, sl.tableData())
}

func (sl SprintList) data() []tui.PreviewData {
	data := make([]tui.PreviewData, 0, len(sl.Data))

	data = append(data, tui.PreviewData{
		Key:  "help",
		Menu: "?",
		Contents: func(s string) interface{} {
			return helpText
		},
	})

	for _, s := range sl.Data {
		bid, sid := s.BoardID, s.ID

		data = append(data, tui.PreviewData{
			Key: fmt.Sprintf("%d-%d-%s", bid, sid, s.StartDate),
			Menu: fmt.Sprintf(
				"➤ #%d %s: ⦗%s - %s⦘",
				s.ID,
				prepareTitle(s.Name),
				formatDateTimeHuman(s.StartDate, time.RFC3339),
				formatDateTimeHuman(s.EndDate, time.RFC3339),
			),
			Contents: func(key string) interface{} {
				if sl.issueCache == nil {
					sl.issueCache = make(map[string]tui.TableData)
				}

				if _, ok := sl.issueCache[key]; !ok {
					issues := sl.Issues(bid, sid)

					sl.issueCache[key] = sl.tabularize(issues)
				}

				return sl.issueCache[key]
			},
		})
	}

	return data
}

func (sl SprintList) tabularize(issues []*jira.Issue) tui.TableData {
	var data tui.TableData

	data = append(data, ValidIssueColumns())

	for _, issue := range issues {
		data = append(data, []string{
			issue.Fields.IssueType.Name,
			issue.Key,
			prepareTitle(issue.Fields.Summary),
			issue.Fields.Status.Name,
			issue.Fields.Assignee.Name,
			issue.Fields.Reporter.Name,
			issue.Fields.Priority.Name,
			issue.Fields.Resolution.Name,
			formatDateTime(issue.Fields.Created, jira.RFC3339),
			formatDateTime(issue.Fields.Updated, jira.RFC3339),
		})
	}

	return data
}

func (sl SprintList) validColumnsMap() map[string]struct{} {
	columns := ValidSprintColumns()
	out := make(map[string]struct{}, len(columns))

	for _, c := range columns {
		out[c] = struct{}{}
	}

	return out
}

func (sl SprintList) tableHeader() []string {
	validColumns, columnsMap := ValidSprintColumns(), sl.validColumnsMap()

	if len(sl.Display.Columns) == 0 {
		return validColumns
	}

	var headers []string

	for _, c := range sl.Display.Columns {
		c = strings.ToUpper(c)

		if _, ok := columnsMap[c]; ok {
			headers = append(headers, strings.ToUpper(c))
		}
	}

	return headers
}

func (sl SprintList) tableData() tui.TableData {
	var data tui.TableData

	headers := sl.tableHeader()

	if !(sl.Display.Plain && sl.Display.NoHeaders) {
		data = append(data, headers)
	}

	if len(headers) == 0 {
		headers = ValidSprintColumns()
	}

	for _, s := range sl.Data {
		data = append(data, sl.assignColumns(headers, s))
	}

	return data
}

func (sl SprintList) assignColumns(columns []string, sprint *jira.Sprint) []string {
	var bucket []string

	for _, column := range columns {
		switch column {
		case fieldID:
			bucket = append(bucket, fmt.Sprintf("%d", sprint.ID))
		case fieldName:
			bucket = append(bucket, sprint.Name)
		case fieldStartDate:
			bucket = append(bucket, formatDateTime(sprint.StartDate, time.RFC3339))
		case fieldEndDate:
			bucket = append(bucket, formatDateTime(sprint.EndDate, time.RFC3339))
		case fieldCompleteDate:
			bucket = append(bucket, formatDateTime(sprint.CompleteDate, time.RFC3339))
		case fieldState:
			bucket = append(bucket, sprint.Status)
		}
	}

	return bucket
}