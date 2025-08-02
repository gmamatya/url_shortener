# URL Shortener

A high-performance URL shortener service built with Go, Fiber, and Redis. This project provides a simple and efficient way to shorten long URLs, with a focus on performance, scalability, and rate limiting.

## Features

- **URL Shortening**: Shorten long URLs to a more manageable length.
- **Custom Aliases**: (Optional) Allow users to specify custom aliases for their shortened URLs.
- **Rate Limiting**: Protect the service from abuse with a robust rate-limiting implementation.
- **High Performance**: Built with Go and the high-performance Fiber web framework.
- **Scalable**: Designed to be scalable and handle a large number of requests.

## Tech Stack

- **Go**: A fast and efficient programming language.
- **Fiber**: A high-performance web framework for Go.
- **Redis**: An in-memory data store used for storing shortened URLs and rate-limiting data.
- **Docker**: A containerization platform for easy deployment and development.

## Getting Started

### Prerequisites

- Go 1.16+
- Docker
- Docker Compose

### Installation

1. Clone the repository:

```bash
git clone https://github.com/gmamatya/url_shortener.git
cd url_shortener
```

2. Create a `.env` file in the `api` directory and add the following environment variables:

```
APP_PORT=3000
DB_ADDR=redis:6379
DB_PASS=
API_QUOTA=10
```

### Running the Service

#### Using Docker Compose

```bash
docker-compose up -d
```

This will start the Go API server and the Redis database in detached mode.

#### Using `go run`

```bash
cd api
go run main.go
```

## API Endpoints

### `POST /api/v1`

Shorten a URL.

**Request Body:**

```json
{
  "url": "https://www.example.com/a-very-long-url",
  "custom_short": "my-custom-alias" 
}
```

**Response:**

```json
{
  "short_url": "http://localhost:3000/my-custom-alias"
}
```

### `GET /{url}`

Resolve a shortened URL.

**Response:**

Redirects to the original URL with a `301 Moved Permanently` status code.

## Rate Limiting

The service implements a rate-limiting mechanism to prevent abuse. The rate limit is based on the client's IP address and is stored in a separate Redis database.

The default rate limit is 10 requests per 30 minutes. This can be configured using the `API_QUOTA` environment variable.

When a client exceeds the rate limit, the service will respond with a `429 Too Many Requests` error.

## Database

The service uses two Redis databases:

- **Database 0**: Stores the shortened URLs.
- **Database 1**: Stores the rate-limiting data.

This separation ensures that the rate-limiting data does not interfere with the primary function of the service.

## Running Tests

To run the tests, use the following command:

```bash
go test ./...
```