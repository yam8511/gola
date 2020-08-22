package discover

import "gola/internal/bootstrap"

// Discover 偽服務發現
func Discover(name string) string {
	if bootstrap.GetAppConf().App.Env == "local" {
		return "localhost:50051"
	}

	return "grpc-" + name + ":50051"
}
