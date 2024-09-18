# Rate Limiter Application

## Overview

This application is a simple HTTP server with a rate-limiting middleware. It limits the number of requests a client can make within a specified time window. The rate limiter is configurable and uses Prometheus for monitoring and visualization of metrics.

## Features

- **Rate Limiting**: Limits the number of requests per IP address within a specified time window.
- **Configurable Limits**: Rate limits can be configured via environment variables.
- **IP Address Parsing**: Handles IP addresses correctly, even when behind a proxy.
- **Logging**: Logs rate-limiting events for monitoring.
- **Metrics Collection**: Uses Prometheus to collect and visualize metrics.

## Configuration

The application uses the following environment variables for configuration:

- `RATE_LIMIT`: The maximum number of requests allowed per IP address within the time window (default: 10).
- `RATE_LIMIT_WINDOW`: The time window for rate limiting (default: 1 minute).

## Setup

### Prerequisites

- Go (https://golang.org/dl/)
- Prometheus (https://prometheus.io/download/)

### Installation

1. **Clone the repository**:

    ```sh
    git clone https://github.com/yourusername/rate-limiter.git
    cd rate-limiter
    ```

2. **Install dependencies**:

    ```sh
    go mod tidy
    ```

3. **Set environment variables** (optional):

    ```sh
    export RATE_LIMIT=10
    export RATE_LIMIT_WINDOW=1m
    ```

4. **Run the application**:

    ```sh
    go run .
    ```

### Prometheus Setup

1. **Install Prometheus**:

    ```sh
    brew install prometheus
    ```

2. **Create a Prometheus configuration file (`prometheus.yml`)**:

    ```yaml
    scrape_configs:
      - job_name: 'prometheus'
        static_configs:
          - targets: ['localhost:9090']
      - job_name: 'rate_limiter'
        static_configs:
          - targets: ['localhost:8080']
    ```

3. **Run Prometheus**:

    ```sh
    prometheus --config.file=prometheus.yml
    ```

4. **Access Prometheus Web UI**:

    Open your web browser and navigate to `http://localhost:9090`.

## Usage

Once the application is running, you can make HTTP requests to `http://localhost:8080`. The rate limiter will allow up to the configured number of requests per IP address within the specified time window. If the limit is exceeded, the server will respond with a `429 Too Many Requests` status.

### Example Request

```sh
curl http://localhost:8080

Metrics
Prometheus metrics are exposed at http://localhost:8080/metrics. You can visualize these metrics using the Prometheus web UI.

License
This project is licensed under the MIT License. See the LICENSE file for details.

Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

Acknowledgements
Prometheus
Go

This README provides a comprehensive overview of the application, including its features, configuration, setup instructions, usage examples, and information on how to visualize metrics using Prometheus.