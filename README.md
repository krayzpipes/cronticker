# cronticker
Golang ticker that works with Cron scheduling.

![Codecov](https://img.shields.io/codecov/c/github/krayzpipes/cronticker)

## Import it

```bash
go get github.com/krayzpipes/cronticker/cronticker
```

```go
import "github.com/krayzpipes/cronticker/cronticker"
```

## Usage

Create a new ticker:

```go
ticker, err := NewTicker("TZ=America/New_York 0 0 0 ? * SUN")
```

Check the ticker's channel for the next tick:

```go
tickerTime := <-ticker.C
```

Reset the ticker to a new cron schedule

```go
err := ticker.Reset("0 0 0 ? * MON,TUE,WED")
```

Stop the ticker
```go
ticker.Stop()
```

Use `defer ticker.Stop()` whenever you can to ensure the cleanup of goroutines.

```go
ticker, _ := NewTicker("@daily")
defer ticker.Stop()
```

### Cron Schedule Format
The Cron schedule can be in **Unix** or **Quartz** format. Directives like
'@weekly' or '@daily' can also be parsed as defined in the
package github.com/robfig/cron/v3.

You may add the TimeZone/location to the beginning of the cron schedule
to change the time zone. Default is UTC if `TZ=Whatever` is not prepended
to the cron schedule.

#### Examples
| Cron Schedule | Type/Format | Description |
|---------------|------|-------------|
|"TZ=America/Los_Angeles 0 0 * * *"|Unix|Daily at 12 AM in Los Angeles|
|"TZ=America/Los_Angeles 0 0 0 ? * MON"|Quartz|Mondays at 12 AM in Los Angeles|
|"TZ=America/Los_Angeles @daily"|Directive|Daily at 12 AM in Los Angeles|
|"@daily"|Directive|Daily at 12 AM UTC|

You may also use the following characters:

| Character | Description | Example |
|-----------|-------------|---------|
|-|Range|4-5|
|/|Step|0/5|
|*|Any|*/5|
|,|List|1,4,5|