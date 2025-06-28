# pau-watcher

A Go-based tool that monitors blockchain transactions for specific addresses and sends notifications via Telegram.

## Features

- Real-time transaction monitoring on multiple blockchain networks (Base, Fantom)
- Telegram notifications for buy/sell transactions
- Configurable watch addresses
- Direct links to DexScreener and DEX trading interfaces
- Automatic filtering of USDC transactions
- Customizable polling intervals

## Prerequisites

- Go 1.x or higher
- Telegram Bot Token
- Telegram Chat ID

## Installation

1. Clone the repository:

```

2. Create an `app.env` file in the project root with the following variables:

```env
TELEGRAM_TOKEN=your_telegram_bot_token
CHAT_ID=your_chat_id
```

## Usage

Run the application with the following command:

```bash
go run main.go --chain <chain_name> --address <wallet_address>
```

### Command Line Arguments

- `--chain`: Specify the blockchain network (required)
  - Supported values: `base`, `fantom`
- `--address`: Wallet address to monitor (optional)
  - Default: "0x2433f77f39815849ede7959c7c43d876242cc4bc"

### Example

```bash
go run main.go --chain base --address 0x123...abc
```

## Project Structure

- `main.go`: Application entry point
- `config/`: Configuration management
- `transaction/`: Transaction parsing and processing
- `telegram/`: Telegram bot integration
- `ticker/`: Transaction monitoring scheduler
- `chain/`: Chain-specific configurations

## Telegram Notifications

The bot sends notifications with the following information:
- Token name
- Blockchain network
- Transaction type (Buy/Sell)
- Interactive buttons for:
  - Trading on the configured DEX
  - Viewing on DexScreener

## Configuration

### Supported Chains

The application currently supports:
- Base Network (BaseScan)
- Fantom Network (FTMScan)

### Monitoring Interval

The default monitoring interval is set to 10 seconds. You can modify this in `main.go`:

```go
ticker := ticker.NewTicker(cfg, 10*time.Second)
```

## Error Handling

The application includes robust error handling for:
- Network requests
- Transaction parsing
- Configuration loading
- Telegram API interactions
