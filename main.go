package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Habit struct to store user responses
type Habit struct {
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
}

var habitFile = "habits.json"

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDString := os.Getenv("TELEGRAM_CHAT_ID")
	log.Println("Retrieving env keys")

	if botToken == "" || chatIDString == "" {
		log.Fatal("Missing required environment variables")
	}
	chatID, err := strconv.ParseInt(chatIDString, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot started...")

	// Send daily habit reminder
	sendHabitReminder(bot, chatID)

	// Listen for user responses
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleUserResponse(bot, chatID, update.Message.Text)
		}
	}
}

// Sends a daily habit reminder
func sendHabitReminder(bot *tgbotapi.BotAPI, chatID int64) {
	message := "🚀 *Daily Habit Tracker* 🚀\n\n" +
		"🔹 *Cognitive:* Read 10 mins\n" +
		"🔹 *Productivity:* Plan your day\n" +
		"🔹 *Emotional Mastery:* Journal for 5 mins\n" +
		"🔹 *Physical:* 20 push-ups\n" +
		"🔹 *Growth:* Learn one new thing\n\n" +
		"Reply ✅ if done, ❌ if missed."

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// Handles user responses (✅ or ❌)
func handleUserResponse(bot *tgbotapi.BotAPI, chatID int64, text string) {
	if text == "✅" || text == "❌" {
		logResponse(text == "✅")
		bot.Send(tgbotapi.NewMessage(chatID, "Got it! ✅"))
	}
}

// Logs responses in a JSON file
func logResponse(completed bool) {
	habits := []Habit{}
	today := time.Now().Format("2006-01-02")

	// Read existing data
	data, err := os.ReadFile(habitFile)
	if err == nil {
		json.Unmarshal(data, &habits)
	}

	// Append today's entry
	habits = append(habits, Habit{Date: today, Completed: completed})

	// Write back to file
	newData, _ := json.Marshal(habits)
	os.WriteFile(habitFile, newData, 0644)
}
