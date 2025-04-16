# AI Shopper Common Library

A collection of very opinionated common utilities and components designed for the AI Shopper project components.

## Features

- **Logger**: A structured logging utility built on top of logrus
- **Kafka**: Kafka client utilities and helpers
- **OS**: Operating system related utilities and helpers
- **Elasticsearch**: Elasticsearch client for document indexing and searching

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

### Elasticsearch

```go
import "netherrealmstudio.com/aishoppercommon/elasticsearch"

// Initialize Elasticsearch client
client, err := elasticsearch.NewElasticsearchClient("localhost", "9200")

// Index a document
err = client.IndexDocument(context.Background(), "my-index", map[string]interface{}{
    "title": "Example Document",
    "content": "This is an example document",
})

// Search documents
results, err := client.SearchDocuments(context.Background(), "my-index", map[string]interface{}{
    "query": map[string]interface{}{
        "match": map[string]interface{}{
            "title": "Example",
        },
    },
})
```

## Requirements

- Go 1.24 or higher
- Kafka (for Kafka-related features)
- Elasticsearch 8.x (for Elasticsearch-related features)

## Dependencies

- github.com/confluentinc/confluent-kafka-go/v2
- github.com/sirupsen/logrus
- github.com/stretchr/testify
- github.com/elastic/go-elasticsearch/v8

## License

[Add your license here]

## Contributing

[Add contribution guidelines if applicable]