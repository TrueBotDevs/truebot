package main

import(
    "fmt"
    "log"
    "math/rand"
    "time"
    "strings"
    "github.com/bwmarrin/discordgo"
) 

func getQuote(s *discordgo.Session, msg *discordgo.MessageCreate, comp string){
    qte, err := db.Query("SELECT quote FROM quotes WHERE quote LIKE \"%"+comp+"%\"")
    if err != nil {
		log.Fatal("Query error:", err)
	}
    defer qte.Close()
    
    var quote string
    var quotes [10000]string
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
        getQuote(s,msg," ")
        return
    }
    newIndex = r1.Intn(index)
    //fmt.Println(quotes)
    //fmt.Println(quotes[newIndex])
    s.ChannelMessageSend(msg.ChannelID,quotes[newIndex])
    fmt.Println(quotes[newIndex])
}

//Cakebombs 10/17
func misQuote(s *discordgo.Session, msg *discordgo.MessageCreate, comp string){
    qte, err := db.Query("SELECT quote FROM quotes WHERE quote LIKE \"%"+comp+"%\"")
    if err != nil {
		log.Fatal("Query error:", err)
	}
    defer qte.Close()
    
    var fakeuser string
    var usercount = 16
    var userchoice = 0
    fakeusers := [16]string{"Ed", "Cakebombs", "Kenos", "Oblivion", "TheTrooble", "Trochlis", "Church", "ZachSK", "Kirkq", "Matty", "Twinge", "Slurpee", "Sent", "z1m", "FearfulFerret", "Muffins"}
    var misquote string
    var quote string
    var quotes [10000]string
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
        getQuote(s,msg," ")
        return
    }
    newIndex = r1.Intn(index)
    //fmt.Println(quotes)
    //fmt.Println(quotes[newIndex])
    //s.ChannelMessageSend(msg.ChannelID,quotes[newIndex])

    userchoice = r1.Intn(usercount)
    fakeuser = fakeusers[userchoice]
    misquote = quotes[newIndex]
    misquote = misquote[:strings.LastIndex(misquote, " ")]
    misquote = misquote + " " + fakeuser
    s.ChannelMessageSend(msg.ChannelID,misquote)

    //fmt.Println(quotes[newIndex])
    fmt.Println(misquote)
}



func addQuote(s *discordgo.Session, msg *discordgo.MessageCreate, quote string){
    if quote == " "{
        s.ChannelMessageSend(msg.ChannelID,"Usage: !addquote <quote>")
    }else if quote == "<quote>"{
        s.ChannelMessageSend(msg.ChannelID,"Very funny Church")
    }else if quote[0] == '"' && strings.Contains(quote, "\" - "){
        if strings.Contains(quote, "<@"){
            s.ChannelMessageSend(msg.ChannelID, "Fuck you, don't @ people in quotes")
        }else{
            newItem := "INSERT INTO quotes (quote,submitter) values (?,?)"
            stmt, err := db.Prepare(newItem)
            if err != nil { panic(err) }
            defer stmt.Close()

            _, err2 := stmt.Exec(quote,msg.Author.Username)
            if err2 != nil { panic(err2) }
            s.ChannelMessageSend(msg.ChannelID, "Added your quote to the database:```" + quote + "```")
        }
    }else{
        s.ChannelMessageSend(msg.ChannelID, "Quotes should be in the format:```\"The thing that was said\" - Username```")
    }
}

func removeQuote(s *discordgo.Session, msg *discordgo.MessageCreate, quote string){

}

func init() {
    CmdList["misquote"] = misQuote
    CmdList["quote"] = getQuote
    CmdList["addquote"] = addQuote
}