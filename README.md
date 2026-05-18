# Majority Judgment Bot for Discord

[![MIT](https://img.shields.io/github/license/MieuxVoter/majority-judgment-bot-discord-golang?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/github/v/release/MieuxVoter/majority-judgment-bot-discord-golang?include_prereleases&style=for-the-badge)](https://github.com/MieuxVoter/majority-judment-bot-discord-golang/releases)
[![Discord Chat https://discord.gg/rAAQG9S](https://img.shields.io/discord/705322981102190593.svg?style=for-the-badge)](https://discord.gg/rAAQG9S)

![Logo of the bot, a weighing scales surrounded by colors.](./doc/mjbot_logo_256.png)

> Helps create Majority Judgment polls in Discord.

![A screenshot of the bot in action.](./doc/screen_00.png)


## Feature Wishlist

- [x] Start a poll with `/mj create …`
- [x] Print some miscellaneous help with `/mj help`
- [x] Inform about my status and metrics with `/mj info`
- [ ] Explain briefly how Majority Judgment works `/mj explain`
- [ ] Record feedback from users with `/mj feedback`
- [ ] Rerun a past poll with `/mj rerun`
- [x] Vote on a poll using buttons
- [x] Look at a poll's result using a button
- [x] Publish a poll's result using a button
- [x] Use only _slash_ (`/`) commands
- [x] Be discreet : do not read messages
- [x] Scope polls per guild (privacy!)
- [x] Enforce quotas per guild
- [x] Choose a `grading` (ex: 👍👎) per poll
- [ ] Add a `secrecy` scope to allow public ballots
- [ ] Docker Compose config (optional)
- [ ] Localization
- [ ] Integration with Liberapay
- [ ] Survive — [🤖🗩 Help!](https://liberapay.com/MajorityJudgmentBot/) — ひとりぼっちのよる


![](./doc/screen_01.png)


## Installation

This bot is in the public beta stage.  Join us on [Discord](https://discord.gg/rAAQG9S) and ask around for an invitation !


## Usage

1. Clone this repository.
2. Create `.env.local`, copied from `.env`
   ```shell
   $ cp .env .env.local
   ```
3. Configure your _discord token_ in `.env.local`
   ```shell
   $ vi .env.local
   ```
4. Run
   ```shell
   $ go run src/main.go
   ```
5. Visit the OAuth URL that was printed in the output


## Build

```shell
$ make
$ ./mjbot
```

> `mjbot` is about 17Mio at the moment, which is way too much.
> We use `upx` in `make release` and it shrinks to `4.5Mio` but it's still too big.


## Using docker

Configure the bot _(discord token, database, log level, etc.)_ in `.env.local`, and run `docker compose`: 

```shell
$ cp .env .env.local
$ vi .env.local
$ docker compose up
```

