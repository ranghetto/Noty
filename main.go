package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

// Config : config.yml
type Config struct {
	Token    string   `yaml:"token"`
	Channels []string `yaml:"channels"`
}

const (
	colorGreen string = "\033[1;32m%s\033[0m"
)

var database *sql.DB
var dbError error
var h Handler
var c Command
var cf Config

func main() {

	cf = Config{}

	data, err := ioutil.ReadFile("config.yml")

	err = yaml.Unmarshal([]byte(data), &cf)
	FatalError(err)

	ds, err := discordgo.New("Bot " + cf.Token)
	FatalError(err)

	err = ds.Open()
	FatalError(err)

	database, dbError = sql.Open("sqlite3", "data/todo.db")
	FatalError(err)

	CreateTableIfDoesntExist(database)

	ds.AddHandler(messagesHandler)

	log.Printf(colorGreen, "Server started... Ctrl+C to stop")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	defer database.Close()
	defer ds.Close()

}

func messagesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	h.SetSession(s)
	h.SetChannelID(m.ChannelID)

	if h.ChannelIsValid(cf) {
		// Controllo che il messaggio non sia stato inviato dal bot stesso
		if m.Author.ID == s.State.User.ID {
			return
		}

		c.SetCommand(m.Content)
		c.SetHandler(h)

		// Elimino il messaggio appena ricevuto
		h.DeleteMessage(m.ID)

		c.WaitAndHandleCommands()
	}

}
