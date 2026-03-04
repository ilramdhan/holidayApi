// Example: Using Holiday API Go SDK
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ilramdhan/holidayapi/pkg/client"
)

func main() {
	// Create a new client
	c := client.New("http://localhost:8080/api/v1")

	ctx := context.Background()

	// Example 1: Get holidays for a specific year
	fmt.Println("=== Holidays in 2024 ===")
	holidays, err := c.GetHolidaysByYear(ctx, 2024)
	if err != nil {
		log.Printf("Error getting holidays: %v", err)
	} else {
		for _, h := range holidays {
			fmt.Printf("%s: %s (%s)\n", h.Date, h.Name, h.Type)
		}
	}

	// Example 2: Get holidays for a specific month
	fmt.Println("\n=== Holidays in January 2024 ===")
	holidays, err = c.GetHolidaysByMonth(ctx, 2024, 1)
	if err != nil {
		log.Printf("Error getting holidays: %v", err)
	} else {
		for _, h := range holidays {
			fmt.Printf("%s: %s\n", h.Date, h.Name)
		}
	}

	// Example 3: Get today's holiday
	fmt.Println("\n=== Today's Holiday ===")
	today, err := c.GetTodayHoliday(ctx)
	if err != nil {
		log.Printf("Error getting today's holiday: %v", err)
	} else if today != nil {
		fmt.Printf("Today is a holiday: %s\n", today.Name)
	} else {
		fmt.Println("Today is not a holiday")
	}

	// Example 4: Get upcoming holidays
	fmt.Println("\n=== Upcoming Holidays (next 5) ===")
	upcoming, err := c.GetUpcomingHolidays(ctx, 5)
	if err != nil {
		log.Printf("Error getting upcoming holidays: %v", err)
	} else {
		for _, h := range upcoming {
			fmt.Printf("%s: %s\n", h.Date, h.Name)
		}
	}

	// Example 5: Get holidays with filters
	fmt.Println("\n=== National Holidays in 2024 ===")
	year := 2024
	filtered, err := c.GetHolidays(ctx, &year, nil, "national")
	if err != nil {
		log.Printf("Error getting filtered holidays: %v", err)
	} else {
		for _, h := range filtered {
			fmt.Printf("%s: %s\n", h.Date, h.Name)
		}
	}

	// Example 6: Health check
	fmt.Println("\n=== Health Check ===")
	health, err := c.HealthCheck(ctx)
	if err != nil {
		log.Printf("Error checking health: %v", err)
	} else {
		fmt.Printf("API Status: %v\n", health["status"])
	}

	// Example 7: Using custom HTTP client with timeout
	fmt.Println("\n=== Using Custom HTTP Client ===")
	customClient := client.New(
		"http://localhost:8080/api/v1",
		client.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}),
	)
	
	holidays, err = customClient.GetHolidaysByYear(ctx, 2024)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Retrieved %d holidays\n", len(holidays))
	}
}
