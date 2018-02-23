package main

import (
    "flag"
    "fmt"
    "strings"
    "github.com/fatih/color"
    "github.com/bwmarrin/discordgo"
    sqlite "github.com/mattn/go-sqlite3"
    "database/sql"
    "log"
    "reflect"
)

// Variables used for command line parameters
var (
    Token string
    BotID string
    CmdList = map[string]interface{}{
        "quote" : getQuote,
    }
    AliasList = map[string]interface{}{
        "quote" : getQuote,
    }
    db *sql.DB
    dgSession *discordgo.Session
    hasSession bool = false
    botTestingChannel = "379073357401948162"
)
//hm
func init() {
    flag.StringVar(&Token, "t", "", "Bot Token")
    flag.Parse()
    
    //Connect to the database
    sql.Register("sqlite3_custom", &sqlite.SQLiteDriver{})
    var err error
    db, err = sql.Open("sqlite3_custom", "./config/TrueBot.db")
    if err != nil {
        log.Fatal("Failed to open database:", err)
	}
	//defer db.Close()
    
    //Function mapping
}

func runInterface(fn interface{},s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    v := reflect.ValueOf(fn)
    rarg := make([]reflect.Value, 3)
    rarg[0] = reflect.ValueOf(s)
    rarg[1] = reflect.ValueOf(msg)
    rarg[2] = reflect.ValueOf(arg)
    v.Call(rarg)
}

func main() {    
    // Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
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

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    dgSession.ChannelMessageSend(botTestingChannel, "Bot is now running TODO: Add a changelog here?")
    
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
    if len(msg.Attachments) >= 1{
        fmt.Fprintf(color.Output, "(%.5s) %s: %s\n", channel.Name, cyan(msg.Author.Username), "uploaded a file.")
    }
    
    //Don't parse empty strings
    if len(msg.Content)==0{
        return
    }
    
    //Echo message to console
    fmt.Fprintf(color.Output, "(%.5s) %s: %s\n", channel.Name, cyan(msg.Author.Username), msg.Content)
    
    //Check for commands
    if (msg.Content[:1] == "!"||msg.Content[:1] == "."){
        cmd := strings.Split(msg.Content, " ")[0][1:]
        var arg = " "
        if len(msg.Content) > len(cmd)+1{
            arg = strings.Replace(msg.Content, "!" + cmd + " ", "", 1)
        }
        
        if CmdList[cmd] != nil{
            runInterface(CmdList[cmd],s,msg,arg)
        }else if AliasList[cmd] !=nil{
            runInterface(AliasList[cmd],s,msg,arg)
        }
    }
}