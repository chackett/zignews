package main

// Config defines configuration that comes from envars
type Config struct {
	MongoAddress  string `env:"MG_ADDR" envDefault:"127.0.0.1:27017"`
	MongoDatabase string `env:"MG_DATABASE" envDefault:"zignews"`
	MongoUser     string `env:"MG_USER" envDefault:"root"`
	MongoPass     string `env:"MG_PASS" envDefault:"password"`
	APIAddress    string `env:"API_ADDR" envDefault:":8080"`
	MsgQueueConn  string `env:"MSG_QUEUE" envDefault:"127.0.0.1:4222"`
}
