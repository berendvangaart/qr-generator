# QR Code Generator API

A simple HTTP service written in Go that generates QR codes with optional watermark overlays.

## Features

- Fast QR code generation
- Optional watermark overlay support
- Docker for easy deployment

## Architecture
```
qr-generator/
├── main.go              # Entry point and dependency wiring
├── handlers/            # HTTP layer (controllers)
├── services/            # Business logic
├── models/              # Data structures
└── utils/               # Helper functions
```

**Dependency Flow**: `main.go` → `handlers` → `services` → `utils`

## Quick Start

### Using Docker

```bash
# Build and run with docker-compose
docker-compose up -d

# Or build and run with Docker directly
docker build -t qr-generator .
docker run -p 8080:8080 qr-generator
```

### Local Development

**Prerequisites:**

- Go 1.24 or higher

```bash
# Install dependencies
go mod download

# Run the server
go run main.go

# Or build and run
go build -o qr-generator
./qr-generator
```

The server will start on `http://localhost:8080`

## API Documentation

### Generate QR Code

**Endpoint:** `POST /generate`

**Content-Type:** `multipart/form-data`

**Parameters:**

| Parameter   | Type   | Required | Description                          |
|-------------|--------|----------|--------------------------------------|
| `content`   | string | Yes      | The content to encode in the QR code |
| `size`      | int    | Yes      | Size of the QR code in pixels        |
| `watermark` | file   | No       | PNG image to overlay on the QR code  |

**Response:**

- Success: PNG image (`image/png`)
- Error: JSON error message with HTTP 400

### Examples

#### Generate Simple QR Code

```bash
curl -X POST http://localhost:8080/generate \
  -F "content=https://example.com" \
  -F "size=256" \
  --output qrcode.png
```

#### Generate QR Code with Watermark

```bash
curl -X POST http://localhost:8080/generate \
  -F "content=https://example.com" \
  -F "size=512" \
  -F "watermark=@logo.png" \
  --output qrcode-with-watermark.png
```

## Development

### Project Structure

- **handlers/qr_handler.go**: HTTP request handling and validation
- **services/qr_service.go**: QR code generation and watermarking logic
- **models/qr_code.go**: Data structures and request models
- **utils/image.go**: Image processing utilities (upload, resize)

### Build Commands

```bash
# Run the application
go run main.go

# Build binary
go build -o qr-generator

# Run tests (if any)
go test ./...

# Format code
go fmt ./...

# Update dependencies
go mod tidy
```

### Docker Commands

```bash
# Build Docker image
docker build -t qr-generator .

# Run container
docker run -p 8080:8080 qr-generator

# Using docker-compose
docker-compose up -d
docker-compose logs -f
docker-compose down
```

## Technical Details

### QR Code Generation

- Uses `github.com/skip2/go-qrcode` library
- Error correction level: Medium
- Supports various sizes (recommended: 256, 512, 1024 pixels)

### Watermark Processing

- Watermark is automatically resized to 25% of QR code width
- Centered overlay using `image/draw` with alpha compositing
- Uses Lanczos3 interpolation for high-quality resizing
- Only PNG format supported for watermarks

### Dependencies

- `github.com/skip2/go-qrcode` - QR code generation
- `github.com/nfnt/resize` - Image resizing with quality interpolation

## Security

- Runs as non-root user in Docker container
- No data stored
- Input validation on all parameters
- Maximum upload size: 10MB
