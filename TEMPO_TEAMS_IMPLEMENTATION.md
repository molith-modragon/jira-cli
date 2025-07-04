# Tempo Teams Command Implementation

## Overview
Successfully added a new command `tempo teams` to the Jira CLI tool that retrieves teams from Tempo.

## Implementation Details

### Files Created/Modified

1. **`internal/cmd/tempo/tempo.go`** - Main tempo command
   - Created the base tempo command structure
   - Follows the same pattern as other commands (project, epic, etc.)
   - Adds the teams subcommand

2. **`internal/cmd/tempo/teams/teams.go`** - Teams subcommand implementation
   - Implements the teams retrieval functionality
   - Supports multiple output formats: table (default), json, and plain text
   - Uses the existing Jira client for authentication
   - Follows the established patterns for command structure

3. **`pkg/tempo/client.go`** - Tempo API client
   - Wraps the existing Jira client to make requests to Tempo API
   - Implements the `GetTeams()` method
   - Uses the Tempo API endpoint: `/rest/tempo-teams/1/team`

4. **`pkg/tempo/types.go`** - Data structures
   - Defines the Team, User, and Member structs
   - Matches the expected Tempo API response format

5. **`pkg/tempo/doc.go`** - Package documentation
   - Provides package-level documentation for the tempo client

6. **`internal/cmd/root/root.go`** - Root command registration
   - Added the tempo command to the main CLI

## Usage

The command can be used in the following ways:

```bash
# List teams in table format (default)
jira tempo teams

# List teams in JSON format
jira tempo teams --format json

# List teams in plain text format
jira tempo teams --plain

# Using the alias
jira tempo team
```

## Command Structure

```
jira tempo teams [flags]

Flags:
  --format string   Output format (table, json) (default "table")
  --plain           Display output in plain text format
  -h, --help        help for teams

Aliases:
  team
```

## Features

- **Multiple Output Formats**: Supports table, JSON, and plain text output
- **Consistent UI**: Uses the same tabwriter and pager system as other commands
- **Authentication**: Leverages existing Jira authentication
- **Error Handling**: Proper error handling and user feedback
- **Help Documentation**: Comprehensive help text and usage instructions

## Technical Notes

- The implementation follows the established patterns in the codebase
- Uses the existing Jira client for authentication and HTTP requests
- Leverages the `tui.PagerOut` function for consistent output formatting
- The Tempo API endpoint used is `/rest/tempo-teams/1/team`
- The command is registered in the main CLI and appears in help menus

## Testing

- Command compiles successfully
- Help text is properly displayed
- Command is registered in the CLI help system
- Follows the same patterns as other commands in the codebase

The implementation is ready for use and follows all the established patterns and conventions of the Jira CLI tool.