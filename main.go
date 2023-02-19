package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
	url2 "net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var regexPattern *regexp.Regexp
var links map[string]string

func loadLinks() {
	yamlFile, err := os.ReadFile("links.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &links)
	if err != nil {
		panic(err)
	}
	for key, value := range links {
		_, _ = fmt.Printf("Loaded replacement link for %s: %s\n", key, value)
	}
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	// load links.yml
	loadLinks()

	discord.AddHandler(messageCreate)
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		regexPattern = regexp.MustCompile(`(https?:\/\/)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
		_ = s.UpdateWatchStatus(0, "chat for links")
	})

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	_ = discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot || m.Author.System {
		return
	}
	matchedUrl := regexPattern.FindString(m.Content)

	if matchedUrl != "" {
		// cleaning up the url
		url, err := url2.Parse(matchedUrl)
		if err != nil {
			return
		}
		url.RawQuery = ""
		// Logging for debug purposes
		fmt.Printf("Link detected in guild %s, channel %s, by %s: %s\n", m.GuildID, m.ChannelID, m.Author.Username, matchedUrl)

		// Matching urls
		// TODO: use maps for matching
		if url.Host == "tiktok.com" || url.Host == "vm.tiktok.com" {
			url.Host = "vm.tiktxk.com"
			_ = s.ChannelTyping(m.ChannelID)
			sendLinkMessage(s, m, matchedUrl, *url)
		} else if url.Host == "instagram.com" || url.Host == "www.instagram.com" {
			url.Host = "ddinstagram.com"
			_ = s.ChannelTyping(m.ChannelID)
			sendLinkMessage(s, m, matchedUrl, *url)
		} else if url.Host == "twitter.com" || url.Host == "www.twitter.com" {
			_ = s.ChannelTyping(m.ChannelID)
			if checkForVideo(s, m) {
				url.Host = "fxtwitter.com"
				sendLinkMessage(s, m, matchedUrl, *url)
			}
		}
	}
	if strings.HasPrefix(m.Content, "<@"+s.State.User.ID+">") && len(m.Content) == len("<@"+s.State.User.ID+">") {
		_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Invite",
						Style: discordgo.LinkButton,
						URL:   "https://discord.com/api/oauth2/authorize?client_id=1073362609115516948&permissions=415001570368&scope=bot",
					},
				},
			}},
		})
	}
}
func checkForVideo(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	time.Sleep(5 * time.Second)
	message, err := s.ChannelMessage(m.ChannelID, m.ID)
	if err != nil {
		return false
	}
	if len(message.Embeds) > 0 {
		if message.Embeds[0] != nil {
			return message.Embeds[0].Video != nil
		}
	}
	return false
}
func sendLinkMessage(s *discordgo.Session, m *discordgo.MessageCreate, match string, url url2.URL) {
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
	newMessage := strings.Replace(m.Content, match, url.String(), 1)
	webhook := getOrCreateWebhook(s, m)
	_, err := s.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Content:   newMessage,
		Username:  m.Author.Username,
		AvatarURL: m.Author.AvatarURL(""),
	})
	if err != nil {
		return
	}
}
func getOrCreateWebhook(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.Webhook {
	webhooks, err := s.ChannelWebhooks(m.ChannelID)
	if err != nil {
		_, _ = s.ChannelMessage(m.ChannelID, "Error getting webhooks: "+err.Error())
		return nil
	}
	for _, webhook := range webhooks {
		if webhook.Name == "LinkFixer" {
			return webhook
		}
	}
	webhook, err := s.WebhookCreate(m.ChannelID, "LinkFixer", s.State.User.AvatarURL(""))
	if err != nil {
		return nil
	}
	return webhook
}
