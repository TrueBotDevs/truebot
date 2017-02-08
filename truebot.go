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
    "time"
    "strconv"
    "reflect"
)

// Variables used for command line parameters
var (
    Token string
    BotID string
    CmdList = map[string]interface{}{
        "quote" : getQuote,
    }
    db *sql.DB
)

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

//week, day, hour, minute, second
func parseDate(date string) (string, time.Duration){
    //1 week 5 days 1 hour 2 minutes 1 second
    compString := "weeks days hours minutes seconds"
    lookingForDates := true
    dateArgs := strings.Split(date," ")
    dateIndex := 0
    var parsedDuration time.Duration
    for lookingForDates {
        if dateIndex >= len(dateArgs)-1{
            lookingForDates = false
            break
        }
        timeInt := strings.Split(date," ")[dateIndex:dateIndex+1][0]
        timeStr := strings.Split(date," ")[dateIndex+1:dateIndex+2][0]
        convertedInt, err := strconv.ParseInt(timeInt,10,32); 
        if err != nil{
            lookingForDates = false
            break
        }
        if strings.Contains(compString,timeStr) == false{
            lookingForDates = false
            break
        }
        fmt.Println(timeInt + " " + timeStr)
        dateIndex += 2
        if strings.Contains("seconds",timeStr){
            parsedDuration += time.Duration(convertedInt)*time.Second
        }else if strings.Contains("days",timeStr){
            parsedDuration += time.Duration(convertedInt*24)*time.Hour
        }else if strings.Contains("hours",timeStr){
            parsedDuration += time.Duration(convertedInt)*time.Hour
        }else if strings.Contains("minutes",timeStr){
            parsedDuration += time.Duration(convertedInt)*time.Minute
        }else if strings.Contains("weeks",timeStr){
            parsedDuration += time.Duration(convertedInt*24*7)*time.Hour
        }
    }
    if dateIndex < len(dateArgs) {
        return strings.Join(strings.Split(date," ")[dateIndex:]," "), parsedDuration
    } else {
        return " ", parsedDuration
    }    
}

func main() {    
    leftover, dur := parseDate("1 week 5 days 1 hour 2 minute 1 second Church is gay")
    fmt.Println(time.Now().Add(dur))
    fmt.Println(leftover)
    //Get row count
    cnt, err := db.Query("SELECT count(*) FROM quotes")
    if err != nil {
		log.Fatal("Query error:", err)
	}
    defer cnt.Close()
    var count int
    for cnt.Next(){
        err = cnt.Scan(&count)
        if err != nil {
            log.Fatal("Parse error:", err)
        }
    }
    
    
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

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
	<-make(chan struct{})
	return
}

func messageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {

    cyan := color.New(color.FgCyan).SprintFunc()
    sender := msg.Author
    //channelID := msg.ChannelID
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    guild, _ := s.Guild(guildID)
    
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
    if msg.Content[:1] == "!"{
        cmd := strings.Split(msg.Content, " ")[0][1:]
        var arg = " "
        if len(msg.Content) > len(cmd)+1{
            arg = strings.Replace(msg.Content, "!" + cmd + " ", "", 1)
        }
        
        if CmdList[cmd] != nil{
            runInterface(CmdList[cmd],s,msg,arg)
        }
        
        if cmd == "ping"{
            s.ChannelMessageSend(msg.ChannelID, "Pong!")
        }        
        if cmd == "addquote"{
            if arg == " "{
                s.ChannelMessageSend(msg.ChannelID,"Usage: !addquote <quote>")
            }else if arg == "<quote>"{
                s.ChannelMessageSend(msg.ChannelID,"Very funny Church")
            }else{
                returnedVal := addQuote(arg)
                s.ChannelMessageSend(msg.ChannelID,returnedVal)
            }
        }
        
        if cmd == "copycat"{
            s.ChannelMessageSend(msg.ChannelID,msg.Content)
        }
        
        if cmd == "azorae"{
            var vChannel *discordgo.Channel
            for _, state := range guild.VoiceStates{
                if state.UserID == sender.ID{
                    v := state.ChannelID
                    vChannel, _ = s.Channel(v)
                }
            }
            if vChannel != nil{
                s.ChannelMessageSend(msg.ChannelID,"<@83742858800009216> you have been pinged to " + vChannel.Name)
                del, _ := s.Channel(msg.ChannelID)
                delThis := del.LastMessageID
                s.ChannelMessageDelete(msg.ChannelID, delThis)
            }
        }
    }
}