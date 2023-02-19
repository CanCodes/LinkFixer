# LinkFixer

LinkFixer is a Discord bot written in Go that automatically fixes broken video embeds. The bot replaces the broken links with the correct links, based on the `links.yaml` file provided in the directory and sends the new version using a webhook that almost looks like you.

![example.gif](assets%2Fexample.gif)
## Usage

To use the bot, simply invite it to your server using the following link:
[Invite Link](https://discord.com/api/oauth2/authorize?client_id=1073362609115516948&permissions=415001570368&scope=bot)

Once the bot is added to your server, simply type a message that contains an embedded video link from Instagram or TikTok. The bot will automatically replace the broken link with the correct link, based on the `links.yaml` file provided in the directory and send the new version using a webhook that almost looks like you.

## Self-Hosting

1. Clone the repository
2. Set up a Discord bot and obtain the bot token.
3. Modify the links.yaml file to configure how LinkFixer handles broken embedded links from different platforms.
4. Build the bot using `go build`.
5. Run the bot using `BOT_TOKEN=<token_string> ./linkfixer`
## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)
