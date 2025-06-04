package packages

import "github.com/rabbytesoftware/quiver/logger"

type ArrowsServer struct {
	Packages 	[]string
	Repository 	string
	
	logs 		*logger.Logger
}

func NewArrowsServer(
	Repository string,
) *ArrowsServer {
	return &ArrowsServer{
		Repository: Repository,
		Packages:  	make([]string, 0),
		logs:		logger.NewLogger("ArrowsServer"),
	}
}
