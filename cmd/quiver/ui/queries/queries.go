package queries

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/client"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/matcher"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/models"

	yaml "gopkg.in/yaml.v3"
)

//go:embed queries.yaml
var queriesByte []byte

type QueryService struct {
	Queries []models.Query

	matcher *matcher.Matcher
	client  *client.Client
}

func NewService(baseURL string) *QueryService {
	queryService := &QueryService{
		client:  client.NewClient(baseURL),
		matcher: matcher.NewMatcher(nil),
	}

	queryService.loadQueries()

	return queryService
}

func (s *QueryService) loadQueries() error {
	if err := s.loadFromMemory(); err != nil {
		return fmt.Errorf("failed to load queries: %w", err)
	}

	s.matcher = matcher.NewMatcher(s.Queries)

	for _, query := range s.Queries {
		if err := s.matcher.ValidateQuery(query); err != nil {
			return fmt.Errorf("invalid query configuration: %w", err)
		}
	}

	return nil
}

func (s *QueryService) HandleCommand(ctx context.Context, input string) (string, int, string, error) {
	match, err := s.matcher.Match(input)
	if err != nil {
		return input, 0, "", fmt.Errorf("failed to match command: %w", err)
	}

	if match.REST == nil {
		return input, 0, "", fmt.Errorf("no matching command found")
	}

	request, err := s.matcher.BuildRequest(match)
	if err != nil {
		return input, 0, "", fmt.Errorf("failed to build request: %w", err)
	}

	response, err := s.client.ExecuteRequest(ctx, request)
	if err != nil {
		return input, 0, "", fmt.Errorf("failed to execute request: %w", err)
	}

	// If the response is not successful, return error with HTTP status and response body
	if !response.Success {
		return input, response.StatusCode, response.GetBodyAsString(), fmt.Errorf("HTTP %d", response.StatusCode)
	}

	return response.String(), response.StatusCode, response.GetBodyAsString(), nil
}

func (s *QueryService) IsLoaded() bool {
	return len(s.Queries) > 0
}

func (s *QueryService) GetAvailableCommands() []string {
	return s.matcher.GetAvailableCommands()
}

func (s *QueryService) GetHelpText() string {
	commands := s.GetAvailableCommands()
	if len(commands) == 0 {
		return "No query commands available."
	}

	var result strings.Builder

	for _, cmd := range commands {
		result.WriteString("  ")
		result.WriteString(cmd)
		result.WriteString("\n")
	}

	result.WriteString("\nNote: Commands with ${argN} require arguments to be provided.")

	return result.String()
}

func (q *QueryService) loadFromMemory() error {
	var config models.QueriesConfig
	if err := yaml.Unmarshal(queriesByte, &config); err != nil {
		return fmt.Errorf("failed to parse queries YAML: %w", err)
	}

	q.Queries = config.Queries
	return nil
}
