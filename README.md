# elf
advent of code discord intergration

## Running
1. Ensure you have [Go 1.17 or newer](https://golang.org/doc/install) installed
2. Run `go mod download`
3. Copy `.env.dist` to `.env` and populate it's contents
   - `ELF_DISCORD_TOKEN`: Discord token for the bot user to run as
   - `ELF_ADVENT_OF_CODE_SESSION`: Advent of Code session cookie for the bot user
4. Run `go run ./cmd/elf`

## Manually Registering Guilds

If you are developing (or if the onboarding process still doesn't exist as of you reading this), it is helpful to manually insert a testing guild into the database. This can be done with the `testdata` tool included.

Example:

```sh
go run ./cmd/testdata --channel-id 909857762064871444 --guild-id 514110851016556567 --leaderboard-code 1111111-11111111 --leaderboard-id 0000001
```