package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token             = fmt.Sprintf("Bot %s", os.Getenv("TOKEN"))
	BotName           = fmt.Sprintf("<@%s>", os.Getenv("CLIENT_ID"))
	stopBot           = make(chan bool)
	vcsession         *discordgo.VoiceConnection
	HelloWorld        = "!helloworld"
	ChannelVoiceJoin  = "!vcjoin"
	ChannelVoiceLeave = "!vcleave"
)

func main() {
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

	fmt.Println("Listening...")
	<-stopBot
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	case strings.HasPrefix(m.Content, "!echo"):
		msg := strings.Fields(m.Content)
		sendMessage(s, c, strings.Join(msg[1:], ""))

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, ChannelVoiceJoin)):

		//今いるサーバーのチャンネル情報の一覧を喋らせる処理を書いておきますね
		//guildChannels, _ := s.GuildChannels(c.GuildID)
		//var sendText string
		//for _, a := range guildChannels{
		//sendText += fmt.Sprintf("%vチャンネルの%v(IDは%v)\n", a.Type, a.Name, a.ID)
		//}
		//sendMessage(s, c, sendText) チャンネルの名前、ID、タイプ(通話orテキスト)をBOTが話す

		//VOICE CHANNEL IDには、botを参加させたい通話チャンネルのIDを代入してください
		//コメントアウトされた上記の処理を使うことでチャンネルIDを確認できます
		vcsession, _ = s.ChannelVoiceJoin(c.GuildID, os.Getenv("VOICE_CHANNEL_ID"), false, false)
		vcsession.AddHandler(onVoiceReceived) //音声受信時のイベントハンドラ

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
