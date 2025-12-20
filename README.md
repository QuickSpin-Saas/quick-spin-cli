# QuickSpin CLI (`qspin`)

The official command-line interface for QuickSpin - a managed microservices SaaS platform providing Redis, RabbitMQ, Elasticsearch, PostgreSQL, and MongoDB as fully managed services.

## Overview

QuickSpin CLI allows developers to provision, manage, and monitor managed services from their terminal. Perfect for developers with limited RAM or Docker configuration challenges, QuickSpin provides enterprise-grade services through a shared Kubernetes multi-tenant architecture.

## Features

- **Service Management**: Create, list, delete, and manage Redis, RabbitMQ, PostgreSQL, MongoDB, MySQL, and Elasticsearch instances
- **Authentication**: Secure JWT-based authentication with OAuth support (GitHub, Google)
- **Multi-Organization**: Switch between organizations and manage team members
- **Billing & Usage**: Monitor usage, view invoices, and manage subscription plans
- **AI-Powered Recommendations**: Get intelligent suggestions for service optimization and cost savings
- **GitOps Support**: Deploy services declaratively using YAML configuration files
- **Multiple Output Formats**: View data as tables, JSON, or YAML
- **Shell Completions**: Auto-completion for bash, zsh, fish, and PowerShell

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap quickspin/tap
brew install qspin
```

### Binary Download

Download the latest release for your platform from [GitHub Releases](https://github.com/quickspin/quickspin-cli/releases).

```bash
# macOS (Apple Silicon)
curl -L https://github.com/quickspin/quickspin-cli/releases/latest/download/qspin-darwin-arm64.tar.gz | tar xz
sudo mv qspin /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/quickspin/quickspin-cli/releases/latest/download/qspin-darwin-amd64.tar.gz | tar xz
sudo mv qspin /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/quickspin/quickspin-cli/releases/latest/download/qspin-linux-amd64.tar.gz | tar xz
sudo mv qspin /usr/local/bin/

# Linux (arm64)
curl -L https://github.com/quickspin/quickspin-cli/releases/latest/download/qspin-linux-arm64.tar.gz | tar xz
sudo mv qspin /usr/local/bin/
```

### Go Install

```bash
go install github.com/quickspin/quickspin-cli/cmd/qspin@latest
```

### Docker

```bash
docker pull quickspin/qspin:latest
docker run -it quickspin/qspin qspin --help
```

## Quick Start

### 1. Initialize Configuration

```bash
qspin config init
```

### 2. Login

```bash
qspin auth login
# Or with flags:
qspin auth login --email user@example.com
```

### 3. Create Your First Service

```bash
# Interactive creation
qspin services create

# Or with flags
qspin services create \
  --name my-redis \
  --type redis \
  --tier developer \
  --region us-east-1
```

### 4. Get Connection Credentials

```bash
# Display credentials
qspin services connect my-redis

# Export as environment variables
qspin services connect my-redis --format env
export $(qspin services connect my-redis --format env | xargs)

# Get connection URI
qspin services connect my-redis --format uri
```

## Usage Examples

### Service Management

```bash
# List all services
qspin services list

# Filter by type and status
qspin services list --type redis --status running

# Get service details
qspin services get my-redis

# View service logs
qspin services logs my-redis --follow --tail 100

# Restart a service
qspin services restart my-redis

# Scale a service
qspin services scale my-redis --tier pro

# Delete a service
qspin services delete my-redis
```

### GitOps Deployment

Create a `quickspin.yaml` file:

```yaml
version: "1"
organization: my-company

services:
  - name: cache-primary
    type: redis
    tier: developer
    region: us-east-1
    config:
      maxmemory: 256mb
      maxmemory-policy: allkeys-lru
    labels:
      environment: production
      team: backend

  - name: message-queue
    type: rabbitmq
    tier: developer
    region: us-east-1
```

Deploy:

```bash
qspin deploy apply -f quickspin.yaml

# Dry run first
qspin deploy apply -f quickspin.yaml --dry-run
```

### AI Recommendations

```bash
# Get service optimization recommendations
qspin ai recommend

# Analyze a specific service
qspin ai analyze --service my-redis

# Get cost optimization suggestions
qspin ai optimize --focus cost
```

### Organization Management

```bash
# List organizations
qspin org list

# Switch organization context
qspin org switch my-company

# View team members
qspin org members

# Invite a team member
qspin org invite developer@example.com --role developer
```

### Billing

```bash
# View current usage
qspin billing usage

# View detailed usage by service
qspin billing usage --detailed

# List invoices
qspin billing invoices

# View current plan
qspin billing plan
```

### Configuration

```bash
# View all configuration
qspin config view

# Set default values
qspin config set defaults.region eu-west-1
qspin config set defaults.output json

# Get specific configuration
qspin config get api.url
```

### Output Formats

All list commands support multiple output formats:

```bash
# Table (default)
qspin services list

# JSON
qspin services list --output json

# YAML
qspin services list --output yaml
```

## Configuration

The CLI stores configuration in `~/.quickspin/config.yaml`. You can customize:

- API endpoint
- Default organization
- Default region and tier
- Output format preferences
- Multiple profiles (production, staging, local)

See [configs/quickspin.example.yaml](configs/quickspin.example.yaml) for a complete example.

## Authentication

QuickSpin CLI supports multiple authentication methods:

### Email/Password Login

```bash
qspin auth login
```

### OAuth (GitHub, Google)

```bash
qspin auth login --browser
```

### API Keys (CI/CD)

For automated workflows:

```bash
export QUICKSPIN_TOKEN=your_api_key
qspin services list
```

## Shell Completions

### Bash

```bash
qspin completion bash > /etc/bash_completion.d/qspin
```

### Zsh

```bash
qspin completion zsh > ~/.zsh/completions/_qspin
```

### Fish

```bash
qspin completion fish > ~/.config/fish/completions/qspin.fish
```

### PowerShell

```powershell
qspin completion powershell > qspin.ps1
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `QUICKSPIN_API_URL` | Override API endpoint |
| `QUICKSPIN_TOKEN` | Authentication token (for CI/CD) |
| `QUICKSPIN_ORG` | Default organization ID |
| `QUICKSPIN_CONFIG` | Config file path |
| `QUICKSPIN_DEBUG` | Enable debug logging |
| `QUICKSPIN_NO_COLOR` | Disable colored output |

## CI/CD Integration

### GitHub Actions

```yaml
- name: Deploy QuickSpin Services
  env:
    QUICKSPIN_TOKEN: ${{ secrets.QUICKSPIN_TOKEN }}
  run: |
    qspin deploy apply -f quickspin.yaml
```

### GitLab CI

```yaml
deploy:
  script:
    - qspin deploy apply -f quickspin.yaml
  variables:
    QUICKSPIN_TOKEN: $CI_QUICKSPIN_TOKEN
```

## Available Service Types

| Service | Description | Tiers Available |
|---------|-------------|-----------------|
| **redis** | In-memory data store and cache | starter, developer, pro, enterprise |
| **rabbitmq** | Message broker | starter, developer, pro, enterprise |
| **postgresql** | Relational database | starter, developer, pro, enterprise |
| **mongodb** | Document database | starter, developer, pro, enterprise |
| **mysql** | Relational database | starter, developer, pro, enterprise |
| **elasticsearch** | Search and analytics engine | developer, pro, enterprise |

## Pricing Tiers

| Tier | CPU | Memory | Storage | Use Case |
|------|-----|--------|---------|----------|
| **starter** | 250m | 512Mi | 1Gi | Development, testing |
| **developer** | 500m | 1Gi | 5Gi | Small projects |
| **pro** | 1000m | 2Gi | 10Gi | Production workloads |
| **enterprise** | 4000m | 8Gi | 50Gi | High-traffic applications |

## Support

- Documentation: https://docs.quickspin.dev
- GitHub Issues: https://github.com/quickspin/quickspin-cli/issues
- Community: https://community.quickspin.dev
- Email: support@quickspin.dev

## Development

### Prerequisites

- Go 1.22+
- Make

### Build from Source

```bash
git clone https://github.com/quickspin/quickspin-cli.git
cd quickspin-cli
make build
./qspin --version
```

### Run Tests

```bash
make test
make test-coverage
```

### Local Development

```bash
# Run against local backend
export QUICKSPIN_API_URL=http://localhost:8000
qspin services list
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
