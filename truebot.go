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
    "math/rand"
    "time"
)

// Variables used for command line parameters
var (
	Token string
	BotID string
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
}

func getQuote(comp string) string{
    qte, err := db.Query("SELECT quote FROM quotes WHERE quote LIKE \"%"+comp+"%\"")
    if err != nil {
		log.Fatal("Query error:", err)
	}
    defer qte.Close()
    
    var quote string
    var quotes [255]string
    var index = 0
    var newIndex = 1
    for qte.Next(){
        err = qte.Scan(&quote)
        if err != nil {
            log.Fatal("Parse error:", err)
        }
        quotes[index] = quote
        index++
    }
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    if index == 0{
        quotes[0] = getQuote(" ")
        index++
    }
    newIndex = r1.Intn(index)
    //fmt.Println(quotes)
    //fmt.Println(quotes[newIndex])
    return quotes[newIndex]
}

func addQuote(quote string) string{
    if strings.Contains(quote, "<@"){
        return "Fuck you, don't @ people in quotes"
    }else{
        //fmt.Println("INSERT INTO quotes (quote) values ('"+quote+"')")
        //_, err := db.Exec("INSERT INTO quotes (quote) values ('"+quote+"')")
        //if err != nil {
        //    log.Fatal(err)
        //}
        //return quote + " Added to the database"
        newItem := "INSERT INTO quotes (quote) values (?)"
        stmt, err := db.Prepare(newItem)
        if err != nil { panic(err) }
        defer stmt.Close()

        _, err2 := stmt.Exec(quote)
        if err2 != nil { panic(err2) }
        return quote + " Added to the database"
    }
}

func main() {    
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
    fmt.Println(count)
    
    //Get a quote
    getQuote("ZachSK")
    
    
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
    if msg.Content[:1] == "!" {
        cmd := strings.Split(msg.Content, " ")[0][1:]
        var arg = " "
        if len(msg.Content) > len(cmd)+1{
            arg = strings.Replace(msg.Content, "!" + cmd + " ", "", 1)
        }
        if cmd == "ping"{
            s.ChannelMessageSend(msg.ChannelID, "Pong!")
        }
        
        if cmd == "quote"{
            quote := getQuote(arg)
            s.ChannelMessageSend(msg.ChannelID,quote)
            fmt.Println(quote)
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
            for _, state := range guild.VoiceStates {
                if state.UserID == sender.ID {
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