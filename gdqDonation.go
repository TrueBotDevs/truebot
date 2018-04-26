package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
)

func gdq(s *discordgo.Session, msg *discordgo.MessageCreate, comp string) {
	resp, err := http.Get("http://taskinoz.com/gdq/api/")
	if err != nil {
		log.Fatal("GDQ  error:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("GDQ  error:", err)
	}
	s.ChannelMessageSend(msg.ChannelID, string(body))
	fmt.Println(string(body))
}

func init() {
	CmdList["donation"] = gdq
}
