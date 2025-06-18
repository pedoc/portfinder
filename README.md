# PortFinder

A cross-platform port finder tool for finding processes using specified ports.

## Features

- 🔍 **Fuzzy Port Matching**: Find processes using ports containing the specified number (e.g., "80" matches 80, 8080, 18080, etc.)
- 📊 **Detailed Process Information**:
  - Process name
  - Process ID (PID)
  - Full command line
  - Working directory
  - Protocol type (TCP/UDP)
  - Connection status (LISTEN, ESTABLISHED, TIME_WAIT, etc.)
  - Port numbers with bound IP addresses
  - Start time
- 🌐 **IPv6 Support**: Properly displays IPv6 addresses with bracket notation
- 🖥️ **Cross-platform**: Windows, macOS, Linux
- ⚡ **Fast Response**: Efficient process lookup and information gathering
- 📈 **Statistics**: Shows total number of processes found

## Installation

### From Source

```bash
git clone <repository-url>
cd portfinder
go build -o pf.exe
```

### Pre-built Binaries

Download pre-built binaries for your platform.

## Usage

```bash
pf.exe <port>
```

### Examples

```bash
# Find processes using ports containing "80"
pf.exe 80

# Find processes using ports containing "8080"
pf.exe 8080

# Find processes using ports containing "22" (SSH)
pf.exe 22
```

### Output Example

```
Port 80(0.0.0.0,127.0.0.1,[::]), 443(0.0.0.0,[::])
Process nginx
PID 12345
Command nginx -g daemon off
WorkDirectory /etc/nginx
Protocol TCP
Status LISTEN
Started 2h

Port 8080(127.0.0.1,[::1])
Process node
PID 48291
Command npm run dev
WorkDirectory ~/projects/test
Protocol TCP
Status LISTEN
Started 3h

Total: 2 processes found
```

## Port Matching

The tool uses **fuzzy matching** - it finds all processes using ports that contain the specified number:

- Input "80" matches: 80, 8080, 18080, 8000, etc.
- Input "22" matches: 22, 2200, 12222, etc.
- Input "443" matches: 443, 4433, 14430, etc.

## IP Address Display

- **IPv4 addresses**: Displayed as-is (e.g., 127.0.0.1, 0.0.0.0)
- **IPv6 addresses**: Displayed with brackets (e.g., [::], [::1], [2001:db8::1])
- **Multiple IPs**: If a port is bound to multiple IPs, they are shown separated by commas

## Connection States

The tool shows all connection states, not just LISTEN:
- **LISTEN**: Port is listening for connections
- **ESTABLISHED**: Active connection
- **TIME_WAIT**: Connection waiting to close
- **CLOSE_WAIT**: Connection waiting for application to close
- And other TCP/UDP states

## Supported Platforms

- Windows
- macOS (Intel & Apple Silicon)
- Linux

## Dependencies

- Go 1.20+
- github.com/shirou/gopsutil/v3

## Building for Different Platforms

Use the provided build script:

```bash
# Windows
.\build.bat

# Or manually
go build -o pf.exe                    # Windows
GOOS=linux GOARCH=amd64 go build -o pf-linux
GOOS=darwin GOARCH=amd64 go build -o pf-macos
GOOS=darwin GOARCH=arm64 go build -o pf-macos-arm64
```

## License

MIT License 