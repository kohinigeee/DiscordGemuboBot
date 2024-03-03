package gemubo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kohinigeee/DiscordGemuboBot/lib"
)

type Preset struct {
	Name     string
	Template *Template
	Params   map[string]string
}

type GemuboMessage struct {
	Title     string
	Content   string
	GuildId   string
	ChannelID string
	GemuboID  string
	MessageID string
	ImageURL  string
	Author    *discordgo.User
	StartTime *time.Time
}

type SpetialPlaceHolder struct {
	Name        string
	Description string
}

var (
	startTime = SpetialPlaceHolder{
		Name:        "$START_TIME",
		Description: "募集事項の開始時間を`hh:mm`形式で指定します",
	}
	title = SpetialPlaceHolder{
		Name:        "$TITLE",
		Description: "募集メッセージのタイトルを指定します",
	}

	imageURL = SpetialPlaceHolder{
		Name:        "$IMAGE_URL",
		Description: "募集メッセージの画像URLを指定します",
	}
)

func SpetialPlaceHolders() []SpetialPlaceHolder {
	return []SpetialPlaceHolder{startTime, title, imageURL}
}

func IsSpetialPlaceHolder(s string) bool {
	for _, sp := range SpetialPlaceHolders() {
		if s == sp.Name {
			return true
		}
	}
	return false
}

func NewPreset(name string, template *Template, params map[string]string) *Preset {
	return &Preset{
		Name:     name,
		Template: template,
		Params:   params,
	}
}

func parseTime(str string) (*time.Time, error) {
	tokens := strings.Split(str, ":")
	hour, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, err
	}

	minu, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	nowJapan := now.Add(time.Hour * 9)

	targetTime := time.Date(nowJapan.Year(), nowJapan.Month(), nowJapan.Day(), hour, minu, 0, 0, time.UTC)
	if targetTime.Before(nowJapan) {
		targetTime = targetTime.Add(time.Hour * 24)
	}
	targetTime = targetTime.Add(time.Hour * -9)
	return &targetTime, nil
}

func (p *Preset) MakeMessage(addtionalParam map[string]string, channelID string, guildID string, author *discordgo.User) (*GemuboMessage, error) {

	msg := p.Template.Content
	params := make(map[string]string, len(p.Params))

	for pname, value := range p.Params {
		params[pname] = value
	}

	if len(addtionalParam) > 0 {
		for pname, value := range addtionalParam {
			params[pname] = value
		}
	}

	gmsg := &GemuboMessage{
		Content:   "",
		StartTime: nil,
		ChannelID: channelID,
		MessageID: "",
		GemuboID:  lib.GeneRandomID(),
		GuildId:   guildID,
		Author:    author,
		ImageURL:  "",
		Title:     "",
	}

	for pname, value := range params {
		switch pname {
		case startTime.Name:
			if value != "NOW" {
				t, err := parseTime(value)
				if err != nil {
					return nil, fmt.Errorf("error: 時間は\"hh:mm\"で指定してください : %w", err)
				}
				gmsg.StartTime = t
			}
		case title.Name:
			gmsg.Title = value
		case imageURL.Name:
			gmsg.ImageURL = value
			gmsg.ImageURL = strings.TrimPrefix(gmsg.ImageURL, "<")
			gmsg.ImageURL = strings.TrimSuffix(gmsg.ImageURL, ">")
		}

		pstr := pname
		msg = strings.ReplaceAll(msg, pstr, value)
	}

	gmsg.Content = msg
	return gmsg, nil
}

func MakeEmbedBosyuMessage(gmsg *GemuboMessage) *discordgo.MessageEmbed {
	msg := ""
	msg += fmt.Sprintf("ID:%s\n", gmsg.GemuboID)
	msg += "─────────────────────────────\n"

	texts := strings.Split(gmsg.Content, "\n")
	for _, text := range texts {
		if text == "" {
			continue
		}
		msg += "### " + text + "\n"
	}

	embed := &discordgo.MessageEmbed{
		Title:       gmsg.Title,
		Description: msg,
		Color:       0x00ff00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    gmsg.Author.Username,
			IconURL: gmsg.Author.AvatarURL("128"),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: gmsg.ImageURL,
		},
	}

	if embed.Title == "" {
		embed.Title = fmt.Sprintf("%sがゲムボ！", gmsg.Author.Username)
	}

	return embed
}
