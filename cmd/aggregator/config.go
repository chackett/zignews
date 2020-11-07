package main

// Config defines configuration that comes from envars
type Config struct {
	MongoAddress  string `env:"MG_ADDR" envDefault:"127.0.0.1:27017"`
	MongoDatabase string `env:"MG_DATABASE" envDefault:"zignews"`
	MongoUser     string `env:"MG_USER" envDefault:"root"`
	MongoPass     string `env:"MG_PASS" envDefault:"password"`
	DelayJobStart bool   `env:"DELAY_START" envDefault:"false"` // Can be used to delay each job start by a random period of time to prevent load spike.
	MsgQueueConn  string `env:"MSG_QUEUE" envDefault:"127.0.0.1:4222"`
}
