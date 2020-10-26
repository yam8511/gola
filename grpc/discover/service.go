package discover

import "gola/internal/bootstrap"

// Discover 偽服務發現
func Discover(name, svc string) string {
	if bootstrap.GetAppConf().App.Env == "local" {
		return "localhost:50051"
	}

	return svc + "-" + name + ":50051"
}
