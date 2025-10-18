# WhatsApp Bot Example

This is a simple example demonstrating how to use the whatsmeow library to create a WhatsApp bot.

## Prerequisites

- Go 1.24 or higher
- A WhatsApp account
- WhatsApp mobile app installed on your phone

## Setup

1. Navigate to the example directory:
```powershell
cd example
```

2. Download dependencies:
```powershell
go mod download
```

3. Build the example:
```powershell
go build
```

4. Run the example:
```powershell
.\example.exe
```

## First Time Login

When you run the example for the first time, it will:
1. Generate a QR code in the terminal
2. Open WhatsApp on your phone
3. Go to Settings → Linked Devices
4. Tap "Link a Device"
5. Scan the QR code displayed in the terminal

After scanning, the session will be saved to `examplestore.db` and you won't need to scan the QR code again.

## What This Example Does

- Connects to WhatsApp using the multidevice API
- Displays received messages in the console
- Shows connection/disconnection events
- Persists the session in a SQLite database

## Modifying the Example

You can modify the event handler in `main.go` to:
- Send automatic replies
- React to specific messages
- Handle different message types (images, videos, documents, etc.)
- Manage groups
- Send messages proactively

## Files Created

- `examplestore.db` - SQLite database storing your session and device information
- `examplestore.db-shm` - SQLite shared memory file
- `examplestore.db-wal` - SQLite write-ahead log

## Cleanup

To start fresh with a new session, delete the database files:
```powershell
Remove-Item examplestore.db*
```

## Documentation

For more information about the whatsmeow API, see:
- [godoc documentation](https://pkg.go.dev/go.mau.fi/whatsmeow)
- [Main repository README](../README.md)
