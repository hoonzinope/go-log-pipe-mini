# go-log-pipe-mini

## Overview

`go-log-pipe-mini` is a high-performance log processing pipeline that monitors log files in real-time, filters logs based on specified conditions, and outputs them to various destinations. Through a flexible configuration file, it allows for easy setup of advanced features such as multiple inputs, multiple outputs, JSON parsing, and log rotation.

## Key Features

-   **Multiple Inputs and Outputs**: Collects logs from multiple sources (`INPUTS`) and sends them to multiple destinations (`OUTPUTS`) simultaneously.
-   **Real-time File Monitoring**: Monitors all log files in a specified folder in real-time to process new logs instantly.
-   **Advanced Log Filtering**:
    -   Supports filtering for plain text (`grep`) and JSON format (`json_grep`) logs.
    -   Allows combining multiple filter conditions with `AND` or `OR` logic.
    -   Supports case-insensitive pattern matching.
-   **Various Output Destinations**:
    -   **Console (`stdout`)**: Displays logs on the standard output.
    -   **File (`file`)**: Saves logs to a file, with support for time-based (`daily`, `hourly`) or size-based (`Rolling`) log rotation.
    -   **HTTP/HTTPS (`http`)**: Sends logs via `POST` request to a specified Webhook URL.
-   **Offset Management**: Saves the processed position of each file in `offset.state` to prevent duplicate processing and ensure data consistency upon restart.
-   **Batch Processing**: Minimizes network and system load by batching multiple logs for transmission with `BATCH_SIZE` and `FLUSH_INTERVAL` options.
-   **Built-in Status Server**: Provides application status (`healthCheck`) and performance metrics (`metrics`) via HTTP endpoints through the `server` module.

## Project Structure

```
/
├───.gitignore
├───config.yml         # Application configuration file
├───go.mod
├───go.sum
├───main.go            # Application entry point
├───README.md
├───confmanager/       # Configuration management
│   └───config.go
├───filter/            # Log filtering
│   ├───filter.go
│   └───manager.go
├───generate/          # Test log generator
│   └───genLog.go
├───input/             # Log input processing
│   ├───manager.go
│   ├───node_file.go
│   ├───node_input.go
│   └───parse/
│       └───json_parser.go
├───offset/            # Offset management
│   └───offset_manager.go
├───output/            # Result output processing
│   ├───console.go
│   ├───file.go
│   ├───httppost.go
│   └───manager.go
├───server/            # Built-in HTTP server
│   ├───healthCheck.go
│   ├───logReciever.go
│   ├───metrics.go
│   └───runServer.go
└───shared/            # Shared data and functions
    ├───data.go
    └───stat.go
```

## Configuration (`config.yml`)

All application behavior is controlled through `config.yml`.

```yaml
INPUTS:
  - NAME: syslog      # Input source identifier
    TYPE: file         # Input type (currently only 'file' is supported)
    PATH: ./logs/      # Path to the folder to monitor
    PARSER: text       # Parser type ('text' or 'json')
  - NAME: applog
    TYPE: file
    PATH: ./json_logs/
    PARSER: json

FILTERS:
  syslog:              # Filters to apply to the 'syslog' input
    MODE: OR           # Filter logic (AND or OR)
    RULES:
      - TYPE: grep     # Plain text filter
        OPTIONS:
          IGNORE_CASE: true
          PATTERN: "error"
  applog:              # Filters to apply to the 'applog' input
    MODE: AND
    RULES:
      - TYPE: json_grep # JSON field filter
        OPTIONS:
          FIELD: "level"
          IGNORE_CASE: true
          PATTERN: "error"

OUTPUTS:
  - TYPE: stdout       # Console output
    TARGETS: [syslog, applog] # Input sources to use for this output
    OPTIONS:
      BATCH_SIZE: 5
      FLUSH_INTERVAL: 2s
  - TYPE: file         # File output
    TARGETS: [syslog]
    OPTIONS:
      PATH: ./output_logs/
      FILENAME: syslog.log
      ROLLING: daily    # Log rotation (daily, hourly, monthly)
      MAX_SIZE: 100MB   # Maximum file size
      MAX_FILES: 7      # Maximum number of retained files
      BATCH_SIZE : 10
      FLUSH_INTERVAL: 5s
  - TYPE: http         # HTTP output
    TARGETS: [syslog, applog]
    OPTIONS:
      URL: http://localhost:8080/logs
      METHOD: POST
      HEADERS:
        Content-Type: application/json
      TIMEOUT: 5s
      BATCH_SIZE: 10
      FLUSH_INTERVAL: 5s
```

## Usage

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/go-log-pipe-mini.git
    cd go-log-pipe-mini
    ```

2.  **Install dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Modify the configuration file**:
    Open `config.yml` and modify `INPUTS`, `FILTERS`, and `OUTPUTS` as needed.

4.  **Run the application**:
    ```bash
    go run main.go
    ```

5.  **Check status (optional)**:
    While the application is running, you can check its status at the following endpoints:
    -   **Health Check**: `http://localhost:8080/health`
    -   **Metrics**: `http://localhost:8080/metrics`

6.  **Generate test logs (optional)**:
    The `generate` package supports creating test logs. You can enable the `generate.GenLogWithFolder` and `generate.GenerateJsonLog` goroutines in `main.go` to test the pipeline.
