package main

import (
    "database/sql"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/fatih/color"
    "github.com/go-ini/ini"
    sqlite "github.com/mattn/go-sqlite3"
    "log"
    "reflect"
    "runtime"
    "strings"
)

// Variables used for command line parameters
var (
    discordKey   string
    BotID        string
    msgOnStartup string
    CmdList      = map[string]interface{}{
        "quote": getQuote,
    }
    AliasList = map[string]interface{}{
        "quote": getQuote,
    }
    db                 *sql.DB
    dgSession          *discordgo.Session
    hasSession         = false
    botCommandsChannel string
    botTestingChannel  string
)

func init() {
    cfg, err := ini.Load("./config/truebot.ini")
    if err != nil {
        fmt.Println("Was not able to load Discord API Key - ", err)
    }
    discordKey = cfg.Section("api-keys").Key("discord").String()
    msgOnStartup = cfg.Section("settings").Key("startup-message").String()
    botCommandsChannel = cfg.Section("channels").Key("bot-commands").String()
    botTestingChannel = cfg.Section("channels").Key("bot-testing").String()

    //Connect to the database
    sql.Register("sqlite3_custom", &sqlite.SQLiteDriver{})
    db, err = sql.Open("sqlite3_custom", "./config/TrueBot.db")
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }
}

func runInterface(fn interface{}, s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    v := reflect.ValueOf(fn)
    rarg := make([]reflect.Value, 3)
    rarg[0] = reflect.ValueOf(s)
    rarg[1] = reflect.ValueOf(msg)
    rarg[2] = reflect.ValueOf(arg)
    v.Call(rarg)
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + discordKey)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
    }
    dgSession = dg
    hasSession = true

    // Get the account information.
    u, err := dg.User("@me")
    if err != nil {
        fmt.Println("error obtaining account details,", err)
    }

    // Store the account ID for later use.
    BotID = u.ID

    // Register messageCreate as a callback for the messageCreate events.
    dg.AddHandler(messageCreate)

    // Register guildMemberAdd as a callback for the GuildMemberAdd events
    dg.AddHandler(guildMemberAdd)

    dg.AddHandler(guildRoleDelete)

    // Open the websocket and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    if msgOnStartup == "True" {
        dgSession.ChannelMessageSend(botTestingChannel, "Bot is now running```"+getLatestChangelog()+"```")
    }

    <-make(chan struct{})
    return
}

func messageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
    channel, _ := s.Channel(msg.ChannelID)
    cyan := color.New(color.FgCyan).SprintFunc()
    //channelID := msg.ChannelID

    // Ignore all messages created by the bot itself
    if msg.Author.ID == BotID {
        return
    }

    //Echo that a user uploaded a file
    if len(msg.Attachments) >= 1 {
        fmt.Fprintf(color.Output, "(%.5s) %s: %s\n", channel.Name, cyan(msg.Author.Username), "uploaded a file.")
    }

    //Don't parse empty strings
    if len(msg.Content) == 0 {
        return
    }

    //Echo message to console
    fmt.Fprintf(color.Output, "(%.5s) %s: %s\n", channel.Name, cyan(msg.Author.Username), msg.Content)

    //Check for commands
    if msg.Content[:1] == "!" || msg.Content[:1] == "." {
        cmd := strings.Split(msg.Content, " ")[0][1:]
        var arg = " "
        if len(msg.Content) > len(cmd)+1 {
            arg = strings.Replace(msg.Content, "!"+cmd+" ", "", 1)
            arg = strings.Replace(arg, "."+cmd+" ", "", 1)
        }

        cmd = strings.ToLower(cmd)
        if CmdList[cmd] != nil {
            runInterface(CmdList[cmd], s, msg, arg)
        } else if AliasList[cmd] != nil {
            runInterface(AliasList[cmd], s, msg, arg)
        }
    }
}

func guildMemberAdd(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
    //assigns pleb role
    s.GuildMemberRoleAdd(member.GuildID, member.User.ID, "167509907095027713")
}

func guildRoleDelete(s *discordgo.Session, delete *discordgo.GuildRoleDelete) {
    _, err := db.Exec("DELETE FROM Roles WHERE id = '" + delete.RoleID + "'")
    if err != nil {
        log.Fatal("Exec error GRD:", err)
    }
}
