# Go ERC-4626 Vault Indexer

## Setup

1. **Set environment variables** (in your shell or `.env` file):

   ```sh
   export RPC_URL="https://mainnet.infura.io/v3/YOUR_INFURA_KEY"
   export CONTRACT_ADDRESS="0xYourVaultAddress"
   export PRIVATE_KEY="your_private_key"
   ```

2. **Install dependencies:**

   ```sh
   go mod tidy
   ```

## Run

```sh
go run main.go
```