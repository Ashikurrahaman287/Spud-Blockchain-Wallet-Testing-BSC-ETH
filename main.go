package main

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// Struct to represent user accounts
type UserAccount struct {
	Address    string
	PrivateKey string
}

// Map to store user accounts
var userAccounts map[int64]*UserAccount

// Function to handle /start command
func startHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the BSC Wallet Bot! Use /help to see available commands.")
	bot.Send(msg)
}

// Function to handle /wallet command
func walletHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	if _, ok := userAccounts[int64(userID)]; ok {
		// You would typically interact with Web3 here to get the wallet balance
		// For demonstration purposes, we'll just reply with a message
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your wallet balance: (Placeholder)")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You haven't created a wallet yet. Use /create_wallet to create one.")
		bot.Send(msg)
	}
}

// Function to handle /create_wallet command
func createWalletHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	if _, ok := userAccounts[int64(userID)]; !ok {
		// Generate new Ethereum address and private key
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Println("Error generating private key:", err)
			return
		}
		address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
		// Store the user's account details
		userAccounts[int64(userID)] = &UserAccount{
			Address:    address,
			PrivateKey: hexutil.Encode(crypto.FromECDSA(privateKey)),
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Wallet created!\nAddress: "+address)
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You have already created a wallet.")
		bot.Send(msg)
	}
}

// Function to handle /private_key command
func privateKeyHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	if _, ok := userAccounts[int64(userID)]; ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your private key: "+userAccounts[int64(userID)].PrivateKey)
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You haven't created a wallet yet. Use /create_wallet to create one.")
		bot.Send(msg)
	}
}

// Function to handle /receive command
func receiveHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	user, ok := userAccounts[int64(userID)]
	if ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your wallet address: "+user.Address)
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You haven't created a wallet yet. Use /create_wallet to create one.")
		bot.Send(msg)
	}
}

// Function to handle /help command
func helpHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Available commands:\n"+
			"/start - Start the bot\n"+
			"/help - Show available commands\n"+
			"/wallet - Wallet information\n"+
			"/create_wallet - Create a new wallet\n"+
			"/private_key - Get your private key\n"+
			"/send [address] - Send crypto to another address\n"+
			"/receive - Get your wallet address",
	)
	bot.Send(msg)
}

func main() {
	// Initialize map for user accounts
	userAccounts = make(map[int64]*UserAccount)

	// Create new bot instance
	bot, err := tgbotapi.NewBotAPI("6719700842:AAHi5LD2itHVh2cwHznSUHkdu1gMpJMf8j8")
	if err != nil {
		log.Fatal("Error creating bot:", err)
	}

	// Set up updates configuration
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Get updates from Telegram
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal("Error getting updates:", err)
	}

	// Process updates
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Handle commands
		switch update.Message.Command() {
		case "start":
			startHandler(update, bot)
		case "help":
			helpHandler(update, bot)
		case "wallet":
			walletHandler(update, bot)
		case "create_wallet":
			createWalletHandler(update, bot)
		case "private_key":
			privateKeyHandler(update, bot)
		case "receive":
			receiveHandler(update, bot)
		}
	}
}
