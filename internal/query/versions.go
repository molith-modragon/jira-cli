package query

// Issue is a query type for issue command.
type Version struct {
	Project string
	Flags   FlagParser

	params *VersionParams
}

// NewIssue creates and initializes a new Release type.
func NewVersion(project string, flags FlagParser) (*Version, error) {
	rel := VersionParams{}
	if err := rel.init(flags); err != nil {
		return nil, err
	}
	return &Version{
		Project: project,
		Flags:   flags,
		params:  &rel,
	}, nil
}

// ReleaseParams is issue command parameters.
type VersionParams struct {
	debug bool
}

func (rel *VersionParams) init(flags FlagParser) error {

	return nil
}
