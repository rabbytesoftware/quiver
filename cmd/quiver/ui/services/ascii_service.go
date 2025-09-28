package services

import (
	"fmt"
	"time"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/events"
	"github.com/rabbytesoftware/quiver/internal/core/metadata"
)

// ASCII art for Quiver CLI startup
const quiverASCIIArt = `
██████████████████████████████
██████████████████████████████
██████████████████████████████
████████              ████████
████████              ████████
████████              ████████
████████                      
████████                      
████████                      
████████                      
███████████████       ████████
███████████████       ████████
███████████████       ████████
`

type ASCIIService struct{}

func NewASCIIService() *ASCIIService {
	return &ASCIIService{}
}

func (s *ASCIIService) GetWelcomeLogLine() events.LogLine {
	intro := fmt.Sprintf("%s %s '%s' \n by %s \n",
		metadata.GetName(),
		metadata.GetVersion(),
		metadata.GetVersionCodename(),
		metadata.GetAuthor(),
	)

	ascii := fmt.Sprintf("%s\n%s", quiverASCIIArt, intro)

	return events.LogLine{
		Text:  ascii,
		Level: "info",
		Time:  time.Now(),
	}
}
