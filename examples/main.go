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
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Option 1: Backward compatible - experiments.yaml auto-derived from same directory as FilePath
	client := aurora.NewClient(aurora.NewFetcherStorage(filefetcher.New(filefetcher.Options{
		FilePath:            "parameters.yaml",
		ExperimentsFilePath: "experiments.yaml",
	})), aurora.ClientOptions{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, opts)),
	})

	// Option 2: Explicit paths for parameters and experiments
	// client := aurora.NewClient(aurora.NewFetcherStorage(filefetcher.New(filefetcher.Options{
	//     FilePath:            "config/parameters.yaml",
	//     ExperimentsFilePath: "config/experiments.yaml",
	// })), aurora.ClientOptions{
	//     Logger: slog.New(slog.NewJSONHandler(os.Stdout, opts)),
	// })

	// Option 3: Experiments only (no parameters file)
	// client := aurora.NewClient(aurora.NewFetcherStorage(filefetcher.New(filefetcher.Options{
	//     ExperimentsFilePath: "experiments.yaml",
	// })), aurora.ClientOptions{
	//     Logger: slog.New(slog.NewJSONHandler(os.Stdout, opts)),
	// })

	err := client.Start(context.Background())
	if err != nil {
		log.Fatalf("failed to start client: %v", err)
	}

	fmt.Println("=== Testing parameters without experiment ===")
	for i := 0; i < 10; i++ {
		attribute := aurora.NewAttribute()
		attribute.Set("subscription_plan", "premium")
		attribute.Set("userID", fmt.Sprintf("user_%d", i))
		resolvedValue := client.GetParameter(context.Background(), "numberOfAttempts", attribute)
		fmt.Printf("User %d: %d attempts\n", i, resolvedValue.Int(-1))
	}

	fmt.Println("\n=== Testing experiment parameters ===")
	for i := 0; i < 10; i++ {
		attribute := aurora.NewAttribute()
		attribute.Set("country", "US")
		attribute.Set("userID", fmt.Sprintf("user_%d", i))
		resolvedValue := client.GetParameter(context.Background(), "checkoutButton", attribute)
		if resolvedValue.Matched() {
			fmt.Printf("User %d (experiment): checkoutButton=%v\n", i, resolvedValue.Value())
		} else {
			fmt.Printf("User %d (fallback): checkoutButton=%v\n", i, resolvedValue.String("default"))
		}
		resolvedValue = client.GetParameter(context.Background(), "titleText", attribute)
		if resolvedValue.Matched() {
			fmt.Printf("User %d (experiment): titleText=%v\n", i, resolvedValue.Value())
		} else {
			fmt.Printf("User %d (fallback): titleText=%v\n", i, resolvedValue.String("default"))
		}
	}
}
