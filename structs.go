package main

type Configuration struct {
	Redis     RedisConfig
	Directory string
}

type RedisConfig struct {
	URI      string
	Password string
	Db       int
}

type Stats struct {
	TotalSize      int64
	FileCount      int
	TimeToComplete int64
	LastUpdated    int64
}
