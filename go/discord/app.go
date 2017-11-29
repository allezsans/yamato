package discord

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/allezsans/yamato/go/pubg"
	"github.com/bwmarrin/discordgo"
)

var (
	Token             = fmt.Sprintf("Bot %s", os.Getenv("TOKEN"))
	BotName           = fmt.Sprintf("<@%s>", os.Getenv("CLIENT_ID"))
	vcsession         *discordgo.VoiceConnection
	HelloWorld        = "!helloworld"
	ChannelVoiceJoin  = "!vcjoin"
	ChannelVoiceLeave = "!vcleave"

	PubgAPIKey = os.Getenv("PUBG_API_KEY")
	pubgClient *pubg.API
)

func Start() {
	discord, err := discordgo.New()
	discord.Token = Token
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	// pubg
	pubgClient, err = pubg.New(PubgAPIKey)
	if err != nil {
		fmt.Printf("error")
		log.Fatal(err)
	}

	fmt.Println("Listening...")
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}

	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	if m.Author.Bot {
		return
	}

	switch {
	case strings.HasPrefix(m.Content, "!echo"):
		msg := strings.Fields(m.Content)
		sendMessage(s, c, strings.Join(msg[1:], ""))

	case strings.HasPrefix(m.Content, "!pubg"):
		msg := strings.Fields(m.Content)
		args := msg[1:]
		if len(args) != 4 {
			sendMessage(s, c, "`!pubg [アカウント名] [na/eu/as/oc/sa/sea/krj] [2017-pre1/2017-pre2/2017-pre3/2017-pre4/2017-pre5] [solo/duo/squad/solo-fpp/duo-fpp/squad-fpp]`")
			return
		}

		player, err := pubgClient.GetPlayer(args)
		if err != nil {
			sendMessage(s, c, "`!pubg [アカウント名] [na/eu/as/oc/sa/sea/krj] [2017-pre1/2017-pre2/2017-pre3/2017-pre4/2017-pre5] [solo/duo/squad/solo-fpp/duo-fpp/squad-fpp]`")
			return
		}

		history, err := pubgClient.GetMatchHistory(player)
		if err != nil {
			sendMessage(s, c, "`!pubg [アカウント名] [na/eu/as/oc/sa/sea/krj] [2017-pre1/2017-pre2/2017-pre3/2017-pre4/2017-pre5] [solo/duo/squad/solo-fpp/duo-fpp/squad-fpp]`")
			return
		}
		fmt.Println((*history)[0])

		f := func(s pubg.Stats) bool {
			param := args[1:]
			return s.Region == param[0] && s.Season == param[1] && s.Mode == param[2]
		}
		stats, err := (*player).GetPlayerStatsFilteredBy(f)

		if err != nil {
			sendMessage(s, c, "`!pubg [アカウント名] [na/eu/as/oc/sa/sea/krj] [2017-pre1/2017-pre2/2017-pre3/2017-pre4/2017-pre5] [solo/duo/squad/solo-fpp/duo-fpp/squad-fpp]`")
			return
		}

		overviewEmbed := NewEmbed().
			SetImage(player.Avatar).
			SetTitle("全ステータス").
			AddField("プレイ回数", pubg.SelectLabel(stats, "Rounds Played")).
			AddField("キルレ", pubg.SelectLabel(stats, "K/D Ratio")).
			AddField("ドン勝率", pubg.SelectLabel(stats, "Win %")).
			AddField("ドン勝回数", pubg.SelectLabel(stats, "Wins")).
			AddField("スキルレート", pubg.SelectLabel(stats, "Rating")).
			AddField("ヘッドショット率", pubg.SelectLabel(stats, "Headshot Kill Ratio")).
			AddField("遠距離キル", pubg.SelectLabel(stats, "Longest Kill")).
			InlineAllFields().
			SetColor(0x00ff00).MessageEmbed

		sendEmbed(s, c, overviewEmbed)

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, ChannelVoiceJoin)):
		vcsession, _ = s.ChannelVoiceJoin(c.GuildID, os.Getenv("VOICE_CHANNEL_ID"), false, false)
		vcsession.AddHandler(onVoiceReceived)

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, ChannelVoiceLeave)):
		vcsession.Disconnect()
	}
}

func onVoiceReceived(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	log.Print("onVoiceReceived")
}

func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func sendEmbed(s *discordgo.Session, c *discordgo.Channel, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(c.ID, embed)

	log.Println(">>> " + embed.Title)
	if err != nil {
		log.Println("Error sending embed: ", err)
	}
}
