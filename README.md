# Majority Judgment Bot for Discord

> Helps create Majority Judgment polls in Discord.

![](./doc/screen_00.png)


## Feature Wishlist

- [x] Start a poll with `/mj create …`
- [x] Vote on a poll using buttons
- [x] Look at a poll's result using a button
- [x] Use only `/` commands
- [x] Do not read messages
- [x] Support multiple guilds
- [ ] Quotas per guild
- [ ] Daily Quotas per guild
- [ ] Publish a poll's result using a button
- [ ] Rerun past poll with `/mj rerun`
- [ ] Display metrics with `/mj info`
- [ ] Explain how Majority Judgment works `/mj explain`
- [ ] Allow/Disallow judges, via nickname or roles, per guild
- [ ] Choose a grading (ex: 👍👎) per poll
- [ ] Docker config
- [ ] Trim the database regularly (CRON)
- [ ] Remove `bot` scope (if doable)
- [ ] Pay for itself `/mj love`


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

