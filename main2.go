package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Struct to represent user accounts
type UserAccount struct {
	Address    string
	PrivateKey string
}

// Map to store user accounts
var userAccounts map[int64]*UserAccount

func init() {
	userAccounts = make(map[int64]*UserAccount)
}

// Function to handle /start command
func startHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the SPUD Wallet Bot! Use /help to see available commands.")
	bot.Send(msg)
}

// Function to handle /send command
func sendHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID

	if _, ok := userAccounts[int64(userID)]; !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You haven't created a wallet yet. Use /create_wallet to create one.")
		bot.Send(msg)
		return
	}

	args := strings.Fields(update.Message.CommandArguments())
	if len(args) != 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command. Usage: /send [address] [amount]")
		bot.Send(msg)
		return
	}

	recipient := args[0]
	amountStr := args[1]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid amount.")
		bot.Send(msg)
		return
	}

	// Here you would construct and send the cryptocurrency transaction
	// For demonstration purposes, let's assume the transaction is successful
	transactionHash := "0x123456789abcdef123456789abcdef123456789abcdef123456789abcdef1234"
	transactionDetails := fmt.Sprintf("Transaction sent successfully!\nRecipient: %s\nAmount: %f ETH\nTransaction Hash: %s", recipient, amount, transactionHash)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, transactionDetails)
	bot.Send(msg)
}

// / Function to handle /wallet command
func walletHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	user, ok := userAccounts[int64(userID)]
	if ok {
		// Call the BscScan API to fetch the wallet balance
		walletAddress := user.Address
		walletBalance, err := fetchWalletBalance(walletAddress)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error fetching wallet balance.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Your wallet balance: %f BNB", walletBalance))
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You haven't created a wallet yet. Use /create_wallet to create one.")
		bot.Send(msg)
	}
}

// Function to fetch wallet balance from the BscScan API
func fetchWalletBalance(walletAddress string) (float64, error) {
	// Construct the API URL
	apiUrl := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=balance&address=%s&apikey=MINPWU6K928WSQI1HSVP7QPGMVC6C81FUQ", walletAddress)

	// Make an HTTP GET request to fetch the wallet balance
	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	// Check if the API request was successful
	if response.Status != "1" {
		return 0, fmt.Errorf("API request failed: %s", response.Message)
	}

	// Parse the wallet balance from the API response
	walletBalance, err := strconv.ParseFloat(response.Result, 64)
	if err != nil {
		return 0, err
	}

	// Convert the wallet balance from wei to BNB (if necessary)
	walletBalance /= 1e18

	return walletBalance, nil
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
		case "send":
			sendHandler(update, bot)
		case "receive":
			receiveHandler(update, bot)
		}
	}
}
