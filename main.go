package main;

import (
	"os"
	"flag"
	"syscall"
	"io/ioutil"
	"os/signal"
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot token")
	flag.Parse()
}

var token string

func main() {
	if token == "" {
		bytes, err := ioutil.ReadFile("token.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		token = strings.TrimSpace(string(bytes))
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Running! Interrupt to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, ".convert") || strings.HasPrefix(m.Content, ".c") {
		conversion := strings.TrimPrefix(strings.TrimPrefix(m.Content, ".convert "), ".c ")
		var quant float64;
		var unit1, unit2 string;
		fmt.Sscanf(conversion, "%f %s to %s", &quant, &unit1, &unit2)

		var out, err = ConvertVal(quant, strings.ToLower(strings.TrimSuffix(unit1, "s")), strings.ToLower(strings.TrimSuffix(unit2, "s")))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Conversion error: %v", err))
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%f %s is %f %s", quant, unit1, out, unit2))
	}
}
