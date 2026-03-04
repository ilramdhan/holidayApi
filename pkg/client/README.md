# Holiday API Go SDK

Official Go SDK for Holiday API Indonesia - Production-ready REST API for Indonesian National Holidays & Joint Leave Days.

## Installation

```bash
go get github.com/ilramdhan/holidayapi/pkg/client
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/ilramdhan/holidayapi/pkg/client"
)

func main() {
    // Create client
    c := client.New("https://api.holidayapi.id/v1")
    
    // Get holidays for 2024
    holidays, err := c.GetHolidaysByYear(context.Background(), 2024)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, h := range holidays {
        fmt.Printf("%s: %s\n", h.Date, h.Name)
    }
}
```

## Features

- 🚀 Simple and intuitive API
- ⚡ Context support for timeouts and cancellation
- 🔧 Custom HTTP client support
- 📝 Type-safe responses
- 🎯 Error handling

## API Reference

### Client Creation

```go
// Basic client
c := client.New("https://api.holidayapi.id/v1")

// With custom HTTP client
c := client.New(
    "https://api.holidayapi.id/v1",
    client.WithHTTPClient(&http.Client{
        Timeout: 10 * time.Second,
    }),
)

// With API key (if required)
c := client.New(
    "https://api.holidayapi.id/v1",
    client.WithAPIKey("your-api-key"),
)
```

### Methods

#### GetHolidaysByYear
Get all holidays for a specific year.

```go
holidays, err := c.GetHolidaysByYear(ctx, 2024)
```

#### GetHolidaysByMonth
Get holidays for a specific year and month.

```go
holidays, err := c.GetHolidaysByMonth(ctx, 2024, 1)
```

#### GetTodayHoliday
Get today's holiday if any.

```go
holiday, err := c.GetTodayHoliday(ctx)
if holiday == nil {
    // Today is not a holiday
}
```

#### GetUpcomingHolidays
Get upcoming holidays.

```go
holidays, err := c.GetUpcomingHolidays(ctx, 5) // Get next 5 holidays
```

#### GetHolidays
Get holidays with filters.

```go
year := 2024
holidays, err := c.GetHolidays(ctx, &year, nil, "national")
```

#### HealthCheck
Check API health status.

```go
health, err := c.HealthCheck(ctx)
```

## Error Handling

The SDK returns detailed errors:

```go
holidays, err := c.GetHolidaysByYear(ctx, 2024)
if err != nil {
    if apiErr, ok := err.(*client.APIError); ok {
        // Handle API error
        fmt.Printf("API Error: %s\n", apiErr.Message)
    } else {
        // Handle other errors
        log.Fatal(err)
    }
}
```

## Examples

See the [examples](../../examples/go) directory for complete examples.

## License

MIT License - see [LICENSE](../../LICENSE) for details.
