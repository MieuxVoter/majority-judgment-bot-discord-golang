# Majority Judgment Bot for Discord

> Helps create Majority Judgment polls in Discord.

![](./doc/screen_00.png)


## Feature Wishlist

- [x] Start a poll
- [x] Vote on a poll
- [x] Display a poll's results
- [x] Use only `/` commands 
- [x] Do not read messages
- [x] Support multiple guilds
- [ ] Quotas per guild
- [ ] Rerun past poll
- [ ] Configure the bot per guild
- [ ] Docker config
- [ ] Remove `bot` scope (if doable)
- [ ] Pay for itself


![](./doc/screen_01.png)


## Usage

1. Clone this repository.
2. Configure your discord token in `.env.local`, copied from `.env`
3. Run
   ```
   go run src/main.go
   ```
4. Visit the OAuth URL that was printed in the output


## Built atop Disgord

This leverages the project https://github.com/andersfylling/disgord

