package main

import (
    "github.com/bwmarrin/discordgo"
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"
)

type Changes struct {
    Logs []struct {
        Semver struct {
            Pretty string `json:"pretty"`
            Major  int    `json:"major"`
            Minor  int    `json:"minor"`
            Patch  int    `json:"patch"`
        } `json:"semver"`
        ChangeLog struct {
            Summary      string   `json:"summary"`
            ItemsChanged []string `json:"itemsChanged"`
        } `json:"changeLog"`
    } `json:"logs"`
}

func getLatestChangelog()(string) {
    jsonFile, err := os.Open("config/changelog.json")
    if err != nil {
        fmt.Println(err)
    }

    defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var changes Changes
    json.Unmarshal(byteValue, &changes)
    formattedString := "TrueBot version " + changes.Logs[0].Semver.Pretty
    formattedString += "\n" + changes.Logs[0].ChangeLog.Summary
    for _, item := range changes.Logs[0].ChangeLog.ItemsChanged{
        formattedString += "\n- " + item
    }
    return formattedString
}

func seeChanges(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    log := getLatestChangelog()
    s.ChannelMessageSend(msg.ChannelID, "```" + log + "```")
}

func init() {
    CmdList["changes"] = seeChanges
}
