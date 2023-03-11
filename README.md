# Majority Judgment Bot for Discord

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
- [x] Use only `/` commands
- [x] Do not read messages
- [x] Support multiple guilds
- [x] Quotas per guild
- [ ] Daily Quotas per guild
- [ ] Publish a poll's result using a button
- [x] Rerun past poll with `/mj rerun`
- [ ] Display metrics with `/mj info`
- [ ] Explain how Majority Judgment works `/mj explain`
- [ ] Allow/Disallow judges, via nickname or roles, per guild
- [ ] Choose a grading (ex: 👍👎) per poll
- [ ] Docker config
- [ ] Docker Compose config
- [ ] Trim the database regularly (CRON)
- [ ] Remove `bot` scope (if doable)
- [ ] Pay for itself `/mj love`


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

   $ go build -o mjbot src/main.go
   $ ./mjbot


## Built atop Disgord

This leverages the project https://github.com/andersfylling/disgord

