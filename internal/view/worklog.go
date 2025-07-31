package view

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// PrintWorklogs prints worklogs for an issue
func PrintWorklogs(worklogs []jira.Worklog, plain bool) {
	if plain {
		printWorklogsPlain(os.Stdout, worklogs)
	} else {
		printWorklogsFormatted(os.Stdout, worklogs)
	}
}

func printWorklogsPlain(w io.Writer, worklogs []jira.Worklog) {
	for _, worklog := range worklogs {
		author := worklog.Author.DisplayName
		if author == "" {
			author = worklog.Author.Name
		}
		
		fmt.Fprintf(w, "ID: %s\n", worklog.ID)
		fmt.Fprintf(w, "Author: %s\n", author)
		fmt.Fprintf(w, "Time Spent: %s (%d seconds)\n", worklog.TimeSpent, worklog.TimeSpentSeconds)
		fmt.Fprintf(w, "Started: %s\n", cmdutil.FormatDateTimeHuman(worklog.Started, jira.RFC3339))
		fmt.Fprintf(w, "Created: %s\n", cmdutil.FormatDateTimeHuman(worklog.Created, jira.RFC3339))
		
		if worklog.Updated != worklog.Created {
			fmt.Fprintf(w, "Updated: %s\n", cmdutil.FormatDateTimeHuman(worklog.Updated, jira.RFC3339))
			
			updateAuthor := worklog.UpdateAuthor.DisplayName
			if updateAuthor == "" {
				updateAuthor = worklog.UpdateAuthor.Name
			}
			if updateAuthor != author {
				fmt.Fprintf(w, "Update Author: %s\n", updateAuthor)
			}
		}
		
		if worklog.Comment != "" {
			fmt.Fprintf(w, "Comment: %s\n", worklog.Comment)
		}
		
		fmt.Fprintln(w, "---")
	}
}

func printWorklogsFormatted(w io.Writer, worklogs []jira.Worklog) {
	if len(worklogs) == 0 {
		return
	}

	header := fmt.Sprintf("%s Worklogs", coloredOut(fmt.Sprintf("%d", len(worklogs)), color.FgWhite, color.Bold))
	fmt.Fprintf(w, "\n%s\n", header)
	
	for i, worklog := range worklogs {
		author := worklog.Author.DisplayName
		if author == "" {
			author = worklog.Author.Name
		}
		
		// Format the worklog header
		meta := fmt.Sprintf(
			"\n %s • %s • %s • %s",
			coloredOut(author, color.FgWhite, color.Bold),
			coloredOut(worklog.TimeSpent, color.FgCyan, color.Bold),
			coloredOut(cmdutil.FormatDateTimeHuman(worklog.Started, jira.RFC3339), color.FgWhite, color.Bold),
			coloredOut(fmt.Sprintf("ID: %s", worklog.ID), color.FgGreen),
		)
		
		// Add update information if updated
		if worklog.Updated != worklog.Created {
			updateAuthor := worklog.UpdateAuthor.DisplayName
			if updateAuthor == "" {
				updateAuthor = worklog.UpdateAuthor.Name
			}
			meta += fmt.Sprintf(" • %s", 
				coloredOut(fmt.Sprintf("Updated by %s", updateAuthor), color.FgYellow))
		}
		
		// Add latest marker for the first worklog (most recent)
		if i == 0 {
			meta += fmt.Sprintf(" • %s", coloredOut("Latest worklog", color.FgMagenta, color.Bold))
		}
		
		fmt.Fprintf(w, "%s\n", meta)
		
		// Add comment if present
		if worklog.Comment != "" {
			comment := strings.TrimSpace(worklog.Comment)
			if comment != "" {
				fmt.Fprintf(w, "\n%s\n", comment)
			}
		}
		
		// Add detailed information
		details := fmt.Sprintf(
			"\n    %s: %d seconds | %s: %s | %s: %s",
			gray("Time"),
			worklog.TimeSpentSeconds,
			gray("Created"),
			cmdutil.FormatDateTimeHuman(worklog.Created, jira.RFC3339),
			gray("Issue ID"),
			worklog.IssueID,
		)
		fmt.Fprintf(w, "%s\n", details)
		
		if i < len(worklogs)-1 {
			fmt.Fprintln(w)
		}
	}
}

// PrintWorklogsWithTempo prints worklogs with Tempo custom attributes
func PrintWorklogsWithTempo(worklogs []jira.WorklogWithTempo, plain bool) {
	if plain {
		printWorklogsWithTempoPlain(os.Stdout, worklogs)
	} else {
		printWorklogsWithTempoFormatted(os.Stdout, worklogs)
	}
}

func printWorklogsWithTempoPlain(w io.Writer, worklogs []jira.WorklogWithTempo) {
	for _, worklog := range worklogs {
		author := worklog.Author.DisplayName
		if author == "" {
			author = worklog.Author.Name
		}
		
		fmt.Fprintf(w, "ID: %s\n", worklog.ID)
		fmt.Fprintf(w, "Author: %s\n", author)
		fmt.Fprintf(w, "Time Spent: %s (%d seconds)\n", worklog.TimeSpent, worklog.TimeSpentSeconds)
		
		if worklog.BillableSeconds != nil {
			billableHours := *worklog.BillableSeconds / 3600.0
			fmt.Fprintf(w, "Billable Hours: %.2f hours (%d seconds)\n", billableHours, *worklog.BillableSeconds)
		}
		
		fmt.Fprintf(w, "Started: %s\n", cmdutil.FormatDateTimeHuman(worklog.Started, jira.RFC3339))
		fmt.Fprintf(w, "Created: %s\n", cmdutil.FormatDateTimeHuman(worklog.Created, jira.RFC3339))
		
		if worklog.Updated != worklog.Created {
			updateAuthor := worklog.UpdateAuthor.DisplayName
			if updateAuthor == "" {
				updateAuthor = worklog.UpdateAuthor.Name
			}
			fmt.Fprintf(w, "Updated: %s by %s\n", cmdutil.FormatDateTimeHuman(worklog.Updated, jira.RFC3339), updateAuthor)
		}
		
		if worklog.Comment != "" {
			fmt.Fprintf(w, "Comment: %s\n", worklog.Comment)
		}
		
		// Display Tempo custom attributes
		if len(worklog.TempoAttributes) > 0 {
			fmt.Fprintf(w, "Custom Attributes:\n")
			for _, attr := range worklog.TempoAttributes {
				fmt.Fprintf(w, "  %s: %s\n", attr.Key, attr.Value)
			}
		}
		
		fmt.Fprintf(w, "\n")
	}
}

func printWorklogsWithTempoFormatted(w io.Writer, worklogs []jira.WorklogWithTempo) {
	header := fmt.Sprintf("%s Worklogs", coloredOut(fmt.Sprintf("%d", len(worklogs)), color.FgWhite, color.Bold))
	fmt.Fprintf(w, "\n%s\n", header)
	
	for i, worklog := range worklogs {
		author := worklog.Author.DisplayName
		if author == "" {
			author = worklog.Author.Name
		}
		
		// Format the worklog header with billable hours if available
		timeInfo := worklog.TimeSpent
		if worklog.BillableSeconds != nil {
			billableHours := *worklog.BillableSeconds / 3600.0
			timeInfo = fmt.Sprintf("%s (%.2fh billable)", worklog.TimeSpent, billableHours)
		}
		
		// Format the worklog header
		meta := fmt.Sprintf(
			"\n %s • %s • %s • %s",
			coloredOut(author, color.FgWhite, color.Bold),
			coloredOut(timeInfo, color.FgCyan, color.Bold),
			coloredOut(cmdutil.FormatDateTimeHuman(worklog.Started, jira.RFC3339), color.FgWhite, color.Bold),
			coloredOut(fmt.Sprintf("ID: %s", worklog.ID), color.FgGreen),
		)
		
		// Add update information if updated
		if worklog.Updated != worklog.Created {
			updateAuthor := worklog.UpdateAuthor.DisplayName
			if updateAuthor == "" {
				updateAuthor = worklog.UpdateAuthor.Name
			}
			updateInfo := fmt.Sprintf(" • Updated by %s on %s", 
				updateAuthor, 
				cmdutil.FormatDateTimeHuman(worklog.Updated, jira.RFC3339))
			meta += coloredOut(updateInfo, color.FgMagenta)
		}
		
		fmt.Fprintf(w, "%s\n", meta)
		
		// Display comment if present
		if worklog.Comment != "" {
			comment := strings.TrimSpace(worklog.Comment)
			fmt.Fprintf(w, "\n%s\n", comment)
		}
		
		// Display Tempo custom attributes
		if len(worklog.TempoAttributes) > 0 {
			fmt.Fprintf(w, "\n%s\n", coloredOut("Custom Attributes:", color.FgYellow, color.Bold))
			for _, attr := range worklog.TempoAttributes {
				fmt.Fprintf(w, "  %s: %s\n", 
					coloredOut(attr.Key, color.FgCyan), 
					coloredOut(attr.Value, color.FgWhite))
			}
		}
		
		if i < len(worklogs)-1 {
			fmt.Fprintln(w)
		}
	}
}