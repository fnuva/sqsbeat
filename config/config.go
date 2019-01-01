// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period             time.Duration `config:"period"`
	Credentials        Credentials   `config:"credentials"`
	Region             string        `config:"region"`
	QueueName          string        `config:"queueName"`
	MaxNumberOfMessage int64         `config:"maxNumberOfMessage"`
	VisibilityTimeout  int64         `config:"visibilityTimeout"`
	WaitTimeSeconds    int64         `config:"waitTimeSeconds"`
}

var DefaultConfig = Config{
	Period:             1 * time.Second,
	MaxNumberOfMessage: 1,
	VisibilityTimeout:  20,
	WaitTimeSeconds:    0,
}
