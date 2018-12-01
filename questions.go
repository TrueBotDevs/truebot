package main

import (
    "github.com/bwmarrin/discordgo"
    "strings"
    "time"
)

func addQuestion(s *discordgo.Session, msg *discordgo.MessageCreate, question string) {
    if question == " " {
        s.ChannelMessageSend(msg.ChannelID, "Usage: !addquestion <question>")
    } else if question == "<question>" {
        s.ChannelMessageSend(msg.ChannelID, "Very funny Church")
    } else if question != "error" {
        if strings.Contains(question, "<@") {
            s.ChannelMessageSend(msg.ChannelID, "You may not @ users in question submissions")
        } else {
            newItem := "INSERT INTO jackbox (question,submitter,date) values (?,?,?)"
            stmt, err := db.Prepare(newItem)
            if err != nil {
                panic(err)
            }
            defer stmt.Close()

            _, err2 := stmt.Exec(question, msg.Author.Username, time.Now().Unix())
            if err2 != nil {
                panic(err2)
            }
            s.ChannelMessageSend(msg.ChannelID, "Added your submission to the database:```"+question+"```")
        }
    } else {
        s.ChannelMessageSend(msg.ChannelID, "Error, question not added to DB")
    }
}

func init() {
    CmdList["addquestion"] = addQuestion
}
