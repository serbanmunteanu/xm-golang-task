package main

import "github.com/serbanmunteanu/xm-golang-task/config"

type Server interface {
	Boot(config *config.WebServerConfig)
}
