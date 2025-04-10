# AI Shopper Common Library

A collection of common utilities and components for the AI Shopper project.

## Features

- **Logger**: A structured logging utility built on top of logrus
- **Kafka**: Kafka client utilities and helpers
- **OS**: Operating system related utilities and helpers

## Installation

```bash
go get netherrealmstudio.com/aishoppercommon
```

## Usage

### Logger

```go
import "netherrealmstudio.com/aishoppercommon/logger"

// Initialize logger
logger := logger.NewLogger()

// Use logger
logger.Info("This is an info message")
logger.Error("This is an error message")
```

### Kafka

```go
import "netherrealmstudio.com/aishoppercommon/kafka"

// Initialize Kafka producer
producer, err := kafka.NewProducer(config)
```

### OS Utilities

```go
import "netherrealmstudio.com/aishoppercommon/os"

// Use OS utilities
```

## Requirements

- Go 1.24 or higher
- Kafka (for Kafka-related features)

## Dependencies

- github.com/confluentinc/confluent-kafka-go/v2
- github.com/sirupsen/logrus
- github.com/stretchr/testify

## License

[Add your license here]

## Contributing

[Add contribution guidelines if applicable]