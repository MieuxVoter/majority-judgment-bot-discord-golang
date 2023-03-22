# Majority Judgment Bot for Discord

![](./doc/mjbot_logo_00.png)

> Helps create Majority Judgment polls in Discord.

![](./doc/screen_00.png)


[![MIT](https://img.shields.io/github/license/MieuxVoter/majority-judgment-bot-discord-golang?style=for-the-badge)](LICENSE.md)
[![Release](https://img.shields.io/github/v/release/MieuxVoter/majority-judgment-bot-discord-golang?include_prereleases&style=for-the-badge)](https://github.com/MieuxVoter/majority-judgment-bot-discord-golang/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/MieuxVoter/majority-judgment-bot-discord-golang/go.yml?style=for-the-badge)](https://github.com/MieuxVoter/majority-judgment-bot-discord-golang/actions/workflows/go.yml)
[![Coverage](https://img.shields.io/codecov/c/github/MieuxVoter/majority-judgment-bot-discord-golang?style=for-the-badge&token=FEUB64HRNM)](https://app.codecov.io/gh/MieuxVoter/majority-judgment-bot-discord-golang/)
[![Code Quality](https://img.shields.io/codefactor/grade/github/MieuxVoter/majority-judgment-bot-discord-golang?style=for-the-badge)](https://www.codefactor.io/repository/github/mieuxvoter/majority-judgment-bot-discord-golang)
[![A+](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/mieuxvoter/majority-judgment-bot-discord-golang)
[![Discord Chat https://discord.gg/rAAQG9S](https://img.shields.io/discord/705322981102190593.svg?style=for-the-badge)](https://discord.gg/rAAQG9S)


## Feature Wishlist

- [x] Start a poll with `/mj create …`
- [x] Vote on a poll using buttons
- [x] Look at a poll's result using a button
- [x] Use only _slash_ (`/`) commands
- [x] Be discreet : do not read messages
- [x] Support multiple guilds
- [x] Enforce quotas per guild
- [x] Rerun a past poll with `/mj rerun`
- [x] Inform about my status and metrics with `/mj info`
- [x] Publish a poll's result using a button
- [x] Docker Compose config (optional)
- [x] Choose a `grading` (ex: 👍👎) per poll
- [x] Add a `secrecy` scope to allow public ballots
- [x] Explain how Majority Judgment works `/mj explain`
- [x] Record feedback from users with `/mj feedback`
- [ ] Survive — [🤖🗩 Help!](https://liberapay.com/MajorityJudgmentBot/)


![](./doc/screen_01.png)


## Usage

1. Clone this repository.
2. Create `.env.local`, copied from `.env`
   ```
   $ cp .env .env.local
   ```
3. Configure your _discord token_ in `.env.local`
4. Run
   ```
   $ go run src/main.go
   ```
5. Visit the OAuth URL that was printed in the output


## Build

```
   $ make
   $ ./mjbot
```


## Using docker

Configure the bot _(discord token, database, log level, etc.)_ in `.env.local`, and run `docker compose`: 

```
   $ cp .env .env.local
   $ docker compose up
```


## Dev Notes

- This leverages the excellent `disgord` https://github.com/andersfylling/disgord
- Trying out Go's dependency injection with `di`.
- Using gateways to communicate with Discord.
- Could perhaps deploy a simple http server at some point, for webhooks.
