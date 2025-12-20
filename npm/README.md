# quickspin-cli

Official CLI for QuickSpin - Managed microservices platform providing Redis, RabbitMQ, Elasticsearch, PostgreSQL, and MongoDB as fully managed services.

## Installation

```bash
npm install -g quickspin-cli
```

## Usage

After installation, the `qspin` command will be available globally:

```bash
# Check version
qspin version

# Login
qspin auth login

# Create a service
qspin services create --name my-redis --type redis --tier developer

# List services
qspin services list

# Get help
qspin --help
```

## Quick Start

1. **Login to QuickSpin**
   ```bash
   qspin auth login
   ```

2. **Create your first service**
   ```bash
   qspin services create --name my-redis --type redis --tier developer
   ```

3. **Get connection credentials**
   ```bash
   qspin services connect my-redis
   ```

## Features

- **Service Management**: Create, list, delete, and manage managed services
- **Authentication**: Secure JWT-based authentication with OAuth support
- **Multi-Organization**: Switch between organizations and manage teams
- **Billing & Usage**: Monitor usage and view invoices
- **AI Recommendations**: Get intelligent service optimization suggestions
- **GitOps Support**: Deploy services using YAML configuration files
- **Shell Completions**: Auto-completion for bash, zsh, fish, and PowerShell

## Supported Services

- Redis (in-memory data store)
- RabbitMQ (message broker)
- PostgreSQL (relational database)
- MongoDB (document database)
- MySQL (relational database)
- Elasticsearch (search and analytics)

## Documentation

For full documentation, visit [https://docs.quickspin.dev](https://docs.quickspin.dev)

## Support

- GitHub Issues: [https://github.com/QuickSpin-Saas/quick-spin-cli/issues](https://github.com/QuickSpin-Saas/quick-spin-cli/issues)
- Email: support@quickspin.dev

## License

MIT
