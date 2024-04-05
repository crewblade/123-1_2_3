package main

import (
	"github.com/crewblade/banner-management-service/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
