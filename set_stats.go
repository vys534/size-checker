package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

func GetTreeSize(path string) (int64, int, error) {
	fileCt := 0
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, 0, err
	}
	var total int64
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return 0, 0, err
			}
			fileCt++
			total += info.Size()
		}
	}
	return total, fileCt, nil
}

func SetStats(path string, redisClient *redis.Client) (*Stats, error) {

	var stats Stats

	start := time.Now().UnixNano()
	size, ct, err := GetTreeSize(path)
	finish := time.Now().UnixNano()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stats.TimeToComplete = finish - start
	stats.LastUpdated = time.Now().Unix()

	e := redisClient.Set(ctx, "sc_time_to_complete", stats.TimeToComplete, 0).Err()
	if e != nil {
		return nil, e
	}
	e = redisClient.Set(ctx, "sc_last_updated", stats.LastUpdated, 0).Err()
	if e != nil {
		return nil, e
	}

	if err != nil {
		return nil, err
	}

	stats.FileCount = ct
	stats.TotalSize = size

	e = redisClient.Set(ctx, "sc_total_size", stats.TotalSize, 0).Err()
	if e != nil {
		return nil, e
	}
	e = redisClient.Set(ctx, "sc_file_count", stats.FileCount, 0).Err()
	if e != nil {
		return nil, e
	}

	return &stats, nil
}
