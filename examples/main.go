package main

import (
	"context"
	"fmt"
	"github.com/tuannguyensn2001/aurora-go/core"
	"github.com/tuannguyensn2001/aurora-go/storage/static"
	"log"
	"log/slog"
)

func main() {
	storage := static.NewStorage("parameters.yaml")
	client := core.NewClient(storage, core.ClientOptions{
		Logger: slog.Default().With("aurora"),
	})
	err := client.Start(context.Background())
	if err != nil {
		log.Fatalf("failed to start client: %v", err)
	}

	for i := 0; i < 10; i++ {
		attribute := core.NewAttribute()
		attribute.Set("subscription_plan", "premium")
		attribute.Set("userID", fmt.Sprintf("user_%d", i))
		resolvedValue := client.GetParameter(context.Background(), "numberOfAttempts", attribute)
		fmt.Println(resolvedValue.Int(-1))
	}
}
