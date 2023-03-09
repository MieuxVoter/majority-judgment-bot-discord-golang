# Majority Judgment Bot for Discord


## Feature Wishlist

- [x] Start a poll
- [x] Vote on a poll
- [x] Display a poll's results
- [x] Use only `/` commands 
- [x] Do not read messages
- [x] Support multiple guilds
- [ ] Quotas per guild
- [ ] Clone past poll
- [ ] Configure the bot per guild
- [ ] Remove `bot` scope (if doable)
- [ ] Pay for itself


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

