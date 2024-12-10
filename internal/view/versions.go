package view

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

// ProjectOption is a functional option to wrap project properties.
type VersionOption func(*Version)

// Project is a project view.
type Version struct {
	data   []*jira.Version
	writer io.Writer
	buf    *bytes.Buffer
}

// NewProject initializes a project.
func NewVersion(data []*jira.Version, opts ...VersionOption) *Version {
	p := Version{
		data: data,
		buf:  new(bytes.Buffer),
	}
	p.writer = tabwriter.NewWriter(p.buf, 0, tabWidth, 1, '\t', 0)

	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

// WithProjectWriter sets a writer for the project.
func WithVersionWriter(w io.Writer) VersionOption {
	return func(p *Version) {
		p.writer = w
	}
}

// Render renders the project view.
func (p Version) Render() error {
	p.printHeader()

	for _, d := range p.data {
		if !d.Released {
			fmt.Fprintf(p.writer, "%s\t%s\t%s\t%s\t%s\n", d.ID, d.UserStartDate, d.UserReleaseDate, prepareTitle(d.Name), d.Description)
		}
	}
	if _, ok := p.writer.(*tabwriter.Writer); ok {
		err := p.writer.(*tabwriter.Writer).Flush()
		if err != nil {
			return err
		}
	}

	return tui.PagerOut(p.buf.String())
}

func (p Version) header() []string {
	return []string{
		"ID",
		"START",
		"RELEASE",
		"NAME",
		"DESC",
	}
}

func (p Version) printHeader() {
	headers := p.header()
	end := len(headers) - 1
	for i, h := range headers {
		fmt.Fprintf(p.writer, "%s", h)
		if i != end {
			fmt.Fprintf(p.writer, "\t")
		}
	}
	fmt.Fprintln(p.writer)
}
