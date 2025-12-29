package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/tuannguyensn2001/aurora-go"
)

func main() {
	storage := aurora.NewStaticStorage("examples/parameters.yaml")
	client := aurora.NewClient(storage, aurora.ClientOptions{
		Logger: slog.Default().With("aurora"),
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
