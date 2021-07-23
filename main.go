package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const (
	Version = "0.0.1"
)

func main() {
	log.Printf(">> Size Checker %s", Version)

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	var configuration Configuration

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("No config file set, %s", err)
		} else {
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     configuration.Redis.URI,
		Password: configuration.Redis.Password,
		DB:       configuration.Redis.Db,
	})

	status := redisClient.Ping(ctx).Err()
	if status != nil {
		log.Fatal("Could not ping Redis database: " + status.Error())
	}
	log.Println("Connected to Redis")

	i, err := os.Stat(configuration.Directory)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("%s does not exist in your filesystem.", configuration.Directory)
		} else {
			log.Fatalf("Can't stat the directory in the config: %v", err)
		}
	}
	if i != nil && !i.IsDir() {
		log.Fatalf("Specified storage path is not a directory.")
	}

	s, e := SetStats(configuration.Directory, redisClient)
	if e != nil {
		log.Fatalf("Failed to set stats: %v", err)
	}
	log.Printf("File count: %d, Size: %d, time to finish: %dns", s.FileCount, s.TotalSize, s.TimeToComplete)
	log.Println("Values set in Redis, exiting.")
}
