package matcher

import (
	"fmt"
	"strings"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/models"
)

type Matcher struct {
	queries []models.Query
}

func NewMatcher(queries []models.Query) *Matcher {
	return &Matcher{
		queries: queries,
	}
}

func (m *Matcher) Match(input string) (*models.Query, error) {
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) == 0 {
		return &models.Query{}, nil
	}

	for _, query := range m.queries {
		if result, err := m.matchQuery(query, parts, 0, make(map[string]string)); err == nil && result != nil {
			return result, nil
		}
	}

	return &models.Query{}, nil
}

func (m *Matcher) matchQuery(query models.Query, parts []string, partIndex int, args map[string]string) (*models.Query, error) {
	executeQuery := query

	if partIndex >= len(parts) {
		return nil, fmt.Errorf("invalid part index")
	}

	syntaxParts, _ := m.parseSyntax(query.Syntax)
	
	if partIndex+len(syntaxParts) > len(parts) {
		return nil, fmt.Errorf("invalid part index")
	}

	executeQuery.Variables = make(map[string]string)
	for k, v := range args {
		executeQuery.Variables[k] = v
	}

	for i, syntaxPart := range syntaxParts {
		inputPart := parts[partIndex+i]
		
		if m.isArgument(syntaxPart) {
			argName := m.extractArgName(syntaxPart)
			executeQuery.Variables[argName] = inputPart
		} else {
			if !strings.EqualFold(syntaxPart, inputPart) {
				return nil, fmt.Errorf("command part mismatch: expected '%s', got '%s'", syntaxPart, inputPart)
			}
		}
	}

	nextPartIndex := partIndex + len(syntaxParts)

	if query.REST != nil && nextPartIndex == len(parts) {
		return &executeQuery, nil
	}

	if len(query.Children) > 0 && nextPartIndex < len(parts) {
		for _, child := range query.Children {
			if result, err := m.matchQuery(child, parts, nextPartIndex, executeQuery.Variables); err == nil {
				return result, nil
			}
		}
	}

	if query.REST != nil {
		return &executeQuery, nil
	}

	return nil, fmt.Errorf("invalid part index")
}

func (m *Matcher) parseSyntax(syntax string) ([]string, []string) {
	parts := strings.Fields(syntax)
	var syntaxParts []string
	var args []string

	for _, part := range parts {
		syntaxParts = append(syntaxParts, part)
		if m.isArgument(part) {
			args = append(args, m.extractArgName(part))
		}
	}

	return syntaxParts, args
}

func (m *Matcher) isArgument(part string) bool {
	return strings.HasPrefix(part, "${") && strings.HasSuffix(part, "}")
}

func (m *Matcher) extractArgName(placeholder string) string {
	if !m.isArgument(placeholder) {
		return ""
	}
	return placeholder[2 : len(placeholder)-1] // Remove ${ and }
}

func (m *Matcher) BuildRequest(query *models.Query) (*models.QueryRequest, error) {
	if query.REST == nil {
		return nil, fmt.Errorf("missing REST configuration")
	}

	url := query.REST.URL
	for argName, argValue := range query.Variables {
		placeholder := "${" + argName + "}"
		url = strings.ReplaceAll(url, placeholder, argValue)
	}

	return &models.QueryRequest{
		URL:    url,
		Method: query.REST.Method,
		Args:   query.Variables,
	}, nil
}

func (m *Matcher) GetAvailableCommands() []string {
	var commands []string
	for _, query := range m.queries {
		commands = append(commands, m.getCommandsFromQuery(query, "")...)
	}
	return commands
}

func (m *Matcher) getCommandsFromQuery(query models.Query, prefix string) []string {
	var commands []string
	
	currentCmd := prefix
	if currentCmd != "" {
		currentCmd += " "
	}
	currentCmd += query.Syntax

	if query.REST != nil {
		description := query.Description
		if description == "" {
			description = "No description available"
		}
		commands = append(commands, fmt.Sprintf("%-30s - %s", currentCmd, description))
	}

	for _, child := range query.Children {
		commands = append(commands, m.getCommandsFromQuery(child, currentCmd)...)
	}

	return commands
}

func (m *Matcher) ValidateQuery(query models.Query) error {
	if strings.TrimSpace(query.Syntax) == "" {
		return fmt.Errorf("query syntax cannot be empty")
	}

	if query.REST != nil {
		if query.REST.URL == "" {
			return fmt.Errorf("REST URL cannot be empty for query: %s", query.Syntax)
		}
		if query.REST.Method == "" {
			return fmt.Errorf("REST method cannot be empty for query: %s", query.Syntax)
		}
		
		validMethods := map[string]bool{
			"GET": true, "POST": true, "PUT": true, "DELETE": true,
			"PATCH": true, "HEAD": true, "OPTIONS": true,
		}
		if !validMethods[strings.ToUpper(query.REST.Method)] {
			return fmt.Errorf("invalid HTTP method '%s' for query: %s", query.REST.Method, query.Syntax)
		}
	}

	for _, child := range query.Children {
		if err := m.ValidateQuery(child); err != nil {
			return fmt.Errorf("invalid child query of '%s': %w", query.Syntax, err)
		}
	}

	return nil
}
