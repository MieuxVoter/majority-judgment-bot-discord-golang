# Majority Judgment Bot for Discord

[![MIT](https://img.shields.io/github/license/MieuxVoter/majority-judgment-bot-discord-golang?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/github/v/release/MieuxVoter/majority-judgment-bot-discord-golang?include_prereleases&style=for-the-badge)](https://github.com/MieuxVoter/majority-judment-bot-discord-golang/releases)
[![Discord Chat https://discord.gg/k9YRuZPSZs](https://img.shields.io/discord/705322981102190593.svg?style=for-the-badge)](https://discord.gg/k9YRuZPSZs)

![Logo of the bot, a weighing scales surrounded by colors.](./doc/mjbot_logo_256.png)

> Helps create Majority Judgment polls in Discord.

![A screenshot of the bot in action.](./doc/screen_00.png)


## Feature Wishlist

- [x] Be discreet : do not read messages
- [x] Use _slash_ (`/`) commands
- [x] Start a poll with `/mj create …`
- [x] Print some miscellaneous help with `/mj help`
- [x] Inform about my status and metrics with `/mj info`
- [x] Explain briefly how Majority Judgment works `/mj explain`
- [ ] ~~Record feedback from users with `/mj feedback`~~
- [ ] ~~Rerun a past poll with `/mj rerun`~~
- [x] Vote on a poll using buttons
- [x] Look at a poll's result using a button
- [x] Publish a poll's result using a button
- [x] Scope polls per guild (privacy!)
- [x] Enforce quotas per guild
- [x] Choose a `grading` (ex: 👍👎) per poll
- [ ] Add a `secrecy` scope to allow public ballots
- [x] Docker Compose config (optional)
- [x] Localization
- [ ] Integration with Liberapay
- [ ] Survive — [🤖🗩 Help!](https://liberapay.com/MajorityJudgmentBot/) — ひとりぼっちのよる


## Installation

This bot is in the _public beta_ stage.  Join us on [Discord](https://discord.gg/k9YRuZPSZs) and ask around for an invitation !


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
5. Visit the OAuth URL
   Replace `{{CLIENT_ID}}` by the Application ID that you created in the Discord dev portal.
   ```
   https://discord.com/api/oauth2/authorize?client_id={{CLIENT_ID}}&permissions=51200&scope=bot+applications.commands
   ```
   The bot scope is not mandatory, but it's nice to have the bot show up in the list of connected people.

## Contribute

### Translations

> [!WARNING]
> Discord only allows a specific list of language codes : https://docs.discord.com/developers/reference#locales

[![Weblate Statistics about this project](https://hosted.weblate.org/widget/majority-judgment-bot-discord/287x66-black.png)](https://hosted.weblate.org/engage/majority-judgment-bot-discord)

We're using the amazing _Weblate_ for translations : https://hosted.weblate.org/engage/majority-judgment-bot-discord/

> 💬 You can add a new language or edit existing translations without ever touching any code.

[![Plot of the translations completions by language](https://hosted.weblate.org/widget/majority-judgment-bot-discord/multi-auto.svg)](https://hosted.weblate.org/engage/majority-judgment-bot-discord/)

> 💡 If you see a sub-100% language you're comfortable with, please consider helping with translations.
> You don't have to do everything.  Every little bit helps.  🥜🐜🐜🐜


## Build

```shell
$ make
$ ./mjbot
```

> [!TIP]
> `mjbot` is about 17Mio at the moment, which is way too much.
> We use `upx` in `make release` and it shrinks to `4.5Mio` but it's still quite big.


## Using docker

Configure the bot _(discord token, database, log level, etc.)_ in `.env.local`, and run `docker compose`: 

```shell
$ cp .env .env.local
$ vi .env.local
$ docker compose up
```

