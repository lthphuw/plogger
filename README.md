# PLogger Structured Logger

[![Go Report Card](https://goreportcard.com/badge/github.com/lthphuw/plogger)](https://goreportcard.com/report/github.com/lthphuw/plogger) [![Go CI](https://github.com/lthphuw/plogger/actions/workflows/go.yml/badge.svg)](https://github.com/lthphuw/plogger/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/lthphuw/plogger/graph/badge.svg?token=GQ28QBZFBZ)](https://codecov.io/gh/lthphuw/plogger)

A lightweight logger for Go with structured output and pluggable formatter/writer.

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Features</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#benchmarks">Benchmarks</a></li>
    <!-- <li><a href="#contributing">Contributing</a></li> -->
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <!-- <li><a href="#acknowledgments">Acknowledgments</a></li> -->
  </ol>
</details>

## Features

- High performance.
- Modular design: easily swap `formatters` & `writers`.
- Supports multiple outputs (`console`, `file`, `discard`) for `writer`.
- Benchmark against `zap`, `zerolog`, `logrus`, etc.

## Getting Started

Note that we only supports the two most recent minor versions of Go.

```
go get github.com/lthphuw/plogger@latest
```

## Usage

### Quick start

```go
formatter, _ := plogger.NewJSONFormatter()
writer, _ := plogger.NewConsoleWriter()

// Create logger options using builder pattern
opts := plogger.NewLoggerOptions().
    SetWriter(writer).
    SetFormatter(formatter).
    SetCaller(true)

// Create logger
logger, _ := plogger.NewLogger(opts)

// Log an info message with additional fields
logger.Info(plogger.NewEntry().
    SetMsg("Hi world").
    SetFieldMap(map[string]any{
        "app": "my-service",
        "env": "prod",
    }))
```

Output will look something like this:

```bash
{
  "app": "my-service",
  "env": "prod",
  "file": ".../main/main.go",
  "func": "main.main",
  "level": "info",
  "line": 22,
  "msg": "Hi world",
  "timestamp": "2025-05-22T15:47:36+07:00"
}
```

### Advanced Usage

#### Custom fields

```go

entry := plogger.NewEntry().
    SetMsg("User login").
    AddField("user_id", 123).
    AddField("ip", "192.168.1.1")

logger.Info(entry)
```

#### Use File Writer (Supports Rotation & Safe Concurrency)

```go
fileWriter, _ := plogger.NewFileWriter(
    plogger.NewFileWriterOptions().SetFilename("/var/log/app.log"),
)

logger, _ := plogger.NewLogger(
    plogger.NewLoggerOptions().
        SetWriter(fileWriter).
        SetFormatter(formatter),
)
```

## Benchmarks

The following benchmarks are adapted from the original [zap](https://github.com/uber-go/zap) logging benchmarks.  
We've reused the same test scenarios and environment setup to ensure a fair comparison, with additional entries for `plogger`.

### Log a message:

| Package            | Time (ns/op) | % to `plogger` | Objects Allocated |
| ------------------ | ------------ | -------------- | ----------------- |
| `plogger` (Ours)   | 554.7        | +0%            | 4 allocs/op       |
| `zap`              | 45.91        | -91.7%         | 0 allocs/op       |
| `zap.Check`        | 42.50        | -92.3%         | 0 allocs/op       |
| `zap.CheckSampled` | 28.65        | -94.8%         | 0 allocs/op       |
| `zap.Sugar`        | 56.20        | -89.9%         | 1 allocs/op       |
| `zap.SugarFormat`  | 1533         | +176.3%        | 58 allocs/op      |
| `zerolog`          | 26.43        | -95.2%         | 0 allocs/op       |
| `zerolog.Check`    | 25.78        | -95.4%         | 0 allocs/op       |
| `zerolog.Format`   | 1470         | +165%          | 58 allocs/op      |
| `go-kit`           | 181.3        | -67.3%         | 8 allocs/op       |
| `slog`             | 190.5        | -65.7%         | 0 allocs/op       |
| `slog.LogAttrs`    | 189.1        | -65.9%         | 0 allocs/op       |
| `apex/log`         | 822.1        | +48.2%         | 4 allocs/op       |
| `log15`            | 2040         | +267.6%        | 17 allocs/op      |
| `logrus`           | 1361         | +145.3%        | 20 allocs/op      |
| `stdlib.Println`   | 146.3        | -73.6%         | 1 allocs/op       |
| `stdlib.Printf`    | 1100         | +98.3%         | 57 allocs/op      |

### Log a message and 10 fields:

| Package            | Time (ns/op) | % to plogger | Objects Allocated |
| ------------------ | ------------ | ------------ | ----------------- |
| `plogger` (Ours)   | 3647         | +0%          | 4 allocs/op       |
| `zap`              | 614.6        | -83.1%       | 5 allocs/op       |
| `zap.Check`        | 616.1        | -83.1%       | 5 allocs/op       |
| `zap.CheckSampled` | 94.63        | -97.4%       | 0 allocs/op       |
| `zap.Sugar`        | 1049         | -71.2%       | 10 allocs/op      |
| `zerolog`          | 260.7        | -92.8%       | 1 allocs/op       |
| `zerolog.Check`    | 261.1        | -92.8%       | 1 allocs/op       |
| `go-kit`           | 2214         | -39.3%       | 56 allocs/op      |
| `slog`             | 2198         | -39.7%       | 41 allocs/op      |
| `slog.LogAttrs`    | 2336         | -36.0%       | 40 allocs/op      |
| `apex/log`         | 10595        | +190.5%      | 64 allocs/op      |
| `log15`            | 11873        | +225.7%      | 73 allocs/op      |
| `logrus`           | 12226        | +235.2%      | 84 allocs/op      |

## License

This project is licensed under the MIT License – see the [LICENSE](LICENSE) file for details.

## Contact

Maintained by [Phu](https://github.com/lthphuw).  
For questions, feedback, or support, please contact: <a href="mailto:lthphuw@gmail.com">lthphuw@gmail.com</a>

<p align="right">(<a href="#readme-top">back to top</a>)</p>
