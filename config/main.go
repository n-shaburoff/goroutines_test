package config

import "time"

type Configurations struct {
	Sender SenderConfigurations
}

type SenderConfigurations struct {
	Count time.Duration
}
