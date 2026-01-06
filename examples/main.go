package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	aurora "github.com/tuannguyensn2001/aurora-go"
	filefetcher "github.com/tuannguyensn2001/aurora-go/fetcher/file"
)

func main() {
	// storage := static.NewStorage(static.Options{
	// 	FilePath: "parameters.yaml",
	// })
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	client := aurora.NewClient(aurora.NewStorage(filefetcher.New(filefetcher.Options{
		FilePath: "parameters.yaml",
	})), aurora.ClientOptions{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, opts)),
	})
	err := client.Start(context.Background())
	if err != nil {
		log.Fatalf("failed to start client: %v", err)
	}

	for i := 0; i < 10; i++ {
		attribute := aurora.NewAttribute()
		attribute.Set("subscription_plan", "premium")
		attribute.Set("userID", fmt.Sprintf("user_%d", i))
		resolvedValue := client.GetParameter(context.Background(), "numberOfAttempts", attribute)
		fmt.Println(resolvedValue.Int(-1))
	}
}
