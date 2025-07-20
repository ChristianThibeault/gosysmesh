# gosysmesh

A secure system and process monitoring tool for local and remote servers via SSH.

## Features

- 🖥️ **Local System Monitoring**: CPU, memory, disk usage and process information
- 🌐 **Remote Monitoring**: Monitor multiple remote servers via SSH
- 🎯 **Process Filtering**: Filter by keywords, users, and other criteria
- 🔒 **Security First**: Host key verification, input validation, command injection prevention
- ⚡ **Flexible Usage**: Run once or continuously with configurable intervals
- 🔗 **SSH Jump Hosts**: Support for proxy jump connections

## Quick Start

### Installation

```bash
git clone https://github.com/ChristianThibeault/gosysmesh
cd gosysmesh
go build
```

### Basic Usage

```bash
# Run once (default behavior)
./gosysmesh start --config example-config.yaml

# Run continuously with monitoring interval
./gosysmesh start --loop --config example-config.yaml
```

### Configuration

Copy `example-config.yaml` to `~/.gosysmesh.yaml` or specify with `--config`:

```yaml
interval: "30s"
monitor:
  local:
    process_filters:
      keywords: ["docker", "nginx"]
      users: ["root", "www-data"]
  remote:
    - host: "192.168.1.100"
      user: "admin"
      port: 22
      ssh_key: "~/.ssh/id_rsa"
      process_filters:
        keywords: ["apache", "mysql"]
```

## Security Features

- **SSH Host Key Verification**: Enabled by default (add hosts to `known_hosts`)
- **Input Validation**: All configuration parameters validated
- **Command Injection Prevention**: Whitelisted commands and parameter sanitization
- **Path Traversal Protection**: Safe file path handling

## SSH Setup

For remote monitoring, ensure:

1. SSH key authentication is configured
2. Your public key is in the remote host's `authorized_keys`
3. Remote hosts are added to your `known_hosts` file

```bash
# Add remote host to known_hosts
ssh-keyscan -H your-remote-host >> ~/.ssh/known_hosts
```

## Examples

### Monitor local system once
```bash
./gosysmesh start
```

### Continuous monitoring every 60 seconds
```bash
./gosysmesh start --loop --config my-config.yaml
```

### Monitor specific processes
```yaml
monitor:
  local:
    process_filters:
      keywords: ["postgres", "redis", "nginx"]
      users: ["postgres", "redis", "www-data"]
```

## Output Format

```
[15:04:05] CPU: 15.2% | MEM: 2048/8192 MB | DISK: 45.2/100.0 GB
├── PID 1234  : /usr/bin/nginx -g daemon off;
│   ├── CPU: 2.1%   MEM: 1.5%
│   └── Start: Mon Jan  1 10:00:00   Stat: S   User: www-data

[15:04:05][server1] CPU: 8.5% | MEM: 1024/4096 MB | DISK: 25.1/50.0 GB
└── PID 5678  : /usr/bin/postgres
    ├── CPU: 0.8%   MEM: 12.3%
    └── Start: Sun Dec 31 09:00:00   Stat: S   User: postgres
```

## Development

### Testing
```bash
go test ./...
```

### Building
```bash
go build -o gosysmesh
```
