package configs

var ENV DotEnv

type DotEnv struct {
	App      App
	DataBase Database
	Queue   Queue
	Cache   Cache
}