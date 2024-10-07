# DDPars

DDPars is a concurrent DNS dump parser written in Go. It extracts domain names from DNS dump files within a specified folder and writes the results to an output file (`domains.txt`). The parser processes files asynchronously using multiple threads for optimal performance.

## Features

- **Asynchronous Processing:** Utilizes 10 worker goroutines to process files concurrently for efficient parsing.
- **Domain Filtering:** Removes unwanted characters like `*` and `_` from domain names:
  - Skips domains containing the `*` symbol entirely.
  - Removes all instances of `_` from domain names.
- **Uniqueness:** Writes unique domain names to the output file, avoiding duplicates.
- **Flexible Input:** Parses all DNS dump files within a specified folder.

## Installation

### Prerequisites

- Go (version 1.21.0 or higher)
- Go module: `github.com/miekg/dns`

### Step-by-Step Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/Lercas/DDPars.git
    ```
2. Change directory to the project:
    ```bash
    cd DDPars
    ```
3. Install the required Go package:
    ```bash
    go get github.com/miekg/dns
    ```

## Usage

Run the tool using the following command:

```bash
go run main.go /input_folder
