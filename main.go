package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const prefix = "!"

// Variables used for command line params
var (
	Token     string
	tokenFile string
)

// Parse arguments
func init() {
	// Accept the Discord bot token from the command line
	debugMode := flag.Bool("d", false, "Debug mode")
	flag.Parse()

	if debugMode != nil && *debugMode {
		fmt.Println("Debug mode enabled")
		tokenFile = "debug.env"
	} else {
		tokenFile = "token.env"
	}

	// Load the .env file
	godotenv.Load(tokenFile)
	// if err != nil {
	// fmt.Println("Error loading .env file")
	// // do not panic here. Enviroments like github does not necessarily include a .env file
	// // panic(err)
	// }
	Token = os.Getenv("TOKEN")
	if Token == "" {
		panic("No token found")
	}

}

// Handles any message being created in the guild, parses them,
// and sends them to the commands module.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Exit function if the message was created by a bot
	if m.Author.Bot {
		return
	}

	// Do nothing if prefix is not present
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// Parse the command by trimming the prefix
	parseCommand(s, m, strings.TrimPrefix(m.Content, prefix))

	// Log some details
	fmt.Println("Message received!\nAuthor: " + m.Author.Username + "\nMessage: " + m.Message.Content)
}

func main() {
	// Create a new discord session
	dg, err := discordgo.New("Bot " + Token)
	dg.AddHandler(addHandlers)
	if err != nil {
		fmt.Println("Error creating new discord session, ", err)
		panic(err)
	}

	// Add a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Only cares about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening a connection to discord, ", err)
		panic(err)
	}

	registerCommands(dg)

	defer dg.Close()

	// Listen until signal is received to end.
	fmt.Println("Caramel Bot is running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	removeCommands(dg)
}
