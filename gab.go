package gab

import (
	"fmt"
	"log"
	"strings"

	tapi "github.com/go-telegram-bot-api/telegram-bot-api"
	iapi "gopkg.in/ahmdrz/goinsta.v2"
)

//GAB stands for GiveAwayBot
type GAB struct {
	InstBot *iapi.Instagram
	TelBot  *tapi.BotAPI
	chats   map[int64]string
}

// New Give Away Bot with telegram token
func New(token string) (*GAB, error) {
	bot, err := tapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("Could not create new bot with provided token: %v", err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &GAB{
		TelBot: bot,
	}, nil
}

func (gab *GAB) RunTelBot() {
	u := tapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := gab.TelBot.GetUpdatesChan(u) //fuck those errors

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		data := strings.Split(update.Message.Text, " ")
		msg := tapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "login":
			user := data[1]
			pass := data[2]
			if err := gab.Login(user, pass); err != nil {
				msg.Text = fmt.Sprintf("Could not login: %v", err)
			}
			msg.Text = fmt.Sprintf("Logged in as %s", user)
		case "follow_all_following_of":
			profile := data[1]
			if err := gab.Follow_all_following_of(profile); err != nil {
				msg.Text = fmt.Sprintf("Could not follow all following of %s: %v", profile, err)
			}
			msg.Text = fmt.Sprintf("Followed all follows of %s", profile)
		default:
			msg.Text = "I don't know that command"
		}
		gab.TelBot.Send(msg)
	}
}

//
func (gab *GAB) Login(username string, password string) error {
	inst := iapi.New(username, password)
	err := inst.Login()
	if err != nil {
		return fmt.Errorf("Login failed: %v", err)
	}
	gab.InstBot = inst
	return nil
}

func (gab *GAB) Follow_all_following_of(profile string) error {
	user, err := gab.InstBot.Profiles.ByName(profile)
	if err != nil {
		return fmt.Errorf("Could not find profile %s: %v", profile, err)
	}

	profiles := user.Following()
	if err != nil {
		return fmt.Errorf("Could not get list of following of profile %s: %v", profile, err)
	}

	for profiles.Next() {
		fmt.Println("Next:", profiles.NextID)
		for _, user := range profiles.Users {
			if err := user.Follow(); err != nil {
				return fmt.Errorf("Could not follow profile %s: %v", user.Username, err)
			}
		}
	}

	return nil
}

/*

 */
