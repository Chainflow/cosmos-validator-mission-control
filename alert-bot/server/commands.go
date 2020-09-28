package server

import (
	"cosmos-validator-mission-control/alert-bot/config"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	client "github.com/influxdata/influxdb1-client/v2"
)

// TelegramAlerting
func TelegramAlerting(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	log.Println("len of update.", len(updates))

	msgToSend := ""

	for update := range updates {
		log.Println("Cmng here..", update.Message.Text)
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Text == "/status" {
			msgToSend = GetStatus(cfg, c)
		} else if update.Message.Text == "/node" {
			msgToSend = NodeStatus(cfg, c)
		} else if update.Message.Text == "/peers" {
			msgToSend = GetPeersCountMsg(cfg, c)
		} else if update.Message.Text == "/help" {
			msgToSend = GetHelp()
		} else {
			msgToSend = "Command not found do /help to know about available commands"
		}

		log.Printf("[%s] %s", update.Message.From.UserName, msgToSend)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgToSend)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

// GetHelp returns the msg to show for /help
func GetHelp() string {
	msg := "List of available commands\n /status - returns validator status, voting power, current block height " +
		"and network block height\n /peers - returns number of connected peers\n /node - return status of caught-up\n" +
		"/help - list out the available commands"

	return msg
}

// GetPeersCountMsg returns the no of peers for /peers
func GetPeersCountMsg(cfg *config.Config, c client.Client) string {
	var msg string

	count := GetPeersCount(cfg, c)
	msg = fmt.Sprintf("No of connected peers %s \n", count)

	return msg
}

// NodeStatus returns the node caught up status /node
func NodeStatus(cfg *config.Config, c client.Client) string {
	var status string

	nodeSync := GetNodeSync(cfg, c)
	status = fmt.Sprintf("Your validator node is %s \n", nodeSync)

	return status
}

// GetStatus returns the status messages for /status
func GetStatus(cfg *config.Config, c client.Client) string {
	var status string

	valStatus := GetValStatusFromDB(cfg, c)
	status = fmt.Sprintf("Your validator is currently  %s \n", valStatus)

	valHeight := GetValidatorBlockHeight(cfg, c)
	status = status + fmt.Sprintf("Validator current block height %s \n", valHeight)

	networkHeight := GetNetworkBlock(cfg, c)
	status = status + fmt.Sprintf("Network current block height %s \n", networkHeight)

	votingPower := GetVotingPowerFromDb(cfg, c)
	status = status + fmt.Sprintf("Voting power of your validator is %s \n", votingPower)

	return status
}
