# [CLI] One Piece Mangadownloader

## How to use it?

1. Set up your `.env` -> create `.env` from `.env.example` with your telegram token and chat id
2. Run `go build` to create an `.exe`
3. use `./mangadownloader.exe -list -local -telegram`
    - `-list`: show list of all manga chapter from OPT
    - `-local`: downloads chapters locally
    - `-telegram`: send chapterpages to your Telegram chat

### .env.example

```
# How to get Bot Token? -> https://core.telegram.org/bots#how-do-i-create-a-bot
TELEGRAM_BOT_TOKEN = <YOUR_TELEGRAM_API_TOKEN>

# How to get Chat ID? -> https://stackoverflow.com/a/38388851
# https://api.telegram.org/bot<TOKEN>/getUpdates
TELEGRAM_CHAT_ID=<Your_ChatId>
```

## Features

-   [x] create a cool cli prompt to use different commands:
-   [x] chapterNumber: select chapterNumber directly
-   [x] list: show list of all manga chapter from OPT
-   [x] local: for downloading chapters locally
-   [x] send chapterpages to your Telegram chat

## Demo

![Demo](./demo.gif)
