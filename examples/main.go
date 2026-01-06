package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	aurora "github.com/tuannguyensn2001/aurora-go"
	filefetcher "github.com/tuannguyensn2001/aurora-go/fetcher/file"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	client := aurora.NewClient(aurora.NewStorage(filefetcher.New(filefetcher.Options{
		FilePath: "parameters.yaml",
	})), aurora.ClientOptions{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, opts)),
	})

	// Register a custom operator for string prefix matching
	client.RegisterOperator("startsWith", func(a, b any) bool {
		strA, okA := a.(string)
		strB, okB := b.(string)
		if !okA || !okB {
			return false
		}
		return strings.HasPrefix(strA, strB)
	})

	// Register a custom operator for modulo checking
	client.RegisterOperator("modulo", func(a, b any) bool {
		// a should be divisible by b
		numA, okA := a.(int)
		numB, okB := b.(int)
		if !okA || !okB || numB == 0 {
			return false
		}
		return numA%numB == 0
	})

	err := client.Start(context.Background())
	if err != nil {
		log.Fatalf("failed to start client: %v", err)
	}

	fmt.Println("=== Testing numberOfAttempts parameter ===")
	for i := 0; i < 10; i++ {
		attribute := aurora.NewAttribute()
		attribute.Set("subscription_plan", "premium")
		attribute.Set("userID", fmt.Sprintf("user_%d", i))
		resolvedValue := client.GetParameter(context.Background(), "numberOfAttempts", attribute)
		fmt.Printf("User %d: %d attempts\n", i, resolvedValue.Int(-1))
	}

	fmt.Println("\n=== Testing enableAuth parameter ===")
	for i := 0; i < 5; i++ {
		attribute := aurora.NewAttribute()
		attribute.Set("country", "VN")
		attribute.Set("age", 20+i)
		attribute.Set("userID", fmt.Sprintf("user_%d", i))
		resolvedValue := client.GetParameter(context.Background(), "enableAuth", attribute)
		fmt.Printf("User %d (age %d): auth=%v\n", i, 20+i, resolvedValue.Boolean(false))
	}

	fmt.Println("\n=== Testing custom operators ===")
	// Test with a parameter that uses custom operators (you would need to add this to parameters.yaml)
	attribute := aurora.NewAttribute()
	attribute.Set("email", "admin@example.com")
	attribute.Set("count", 10)
	resolvedValue := client.GetParameter(context.Background(), "customOperatorTest", attribute)
	fmt.Printf("Custom operator test result: %v\n", resolvedValue.Boolean(false))
}
