package main

import(
    "fmt"
    "log"
    "math/rand"
    "time"
    "strings"
    "github.com/bwmarrin/discordgo"
) 


var fakeusers = [16]string{"Ed", "Cakebombs", "Kenos", "Oblivion", "TheTrooble", "Trochlis", "Church", "ZachSK", "Kirkq", "Matty", "Twinge", "Slurpee", "Sent", "z1m", "FearfulFerret", "Muffins"}
var usercount = 16

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
//Trooble sucks virgina
//Cakebombs 10/17
func misQuote(s *discordgo.Session, msg *discordgo.MessageCreate, comp string){
    qte, err := db.Query("SELECT quote FROM quotes WHERE quote LIKE \"%"+comp+"%\"")
    if err != nil {
		log.Fatal("Query error:", err)
	}
    defer qte.Close()
    
    var fakeuser string
    //var usercount = 16
    var userchoice = 0
    //fakeusers := [16]string{"Ed", "Cakebombs", "Kenos", "Oblivion", "TheTrooble", "Trochlis", "Church", "ZachSK", "Kirkq", "Matty", "Twinge", "Slurpee", "Sent", "z1m", "FearfulFerret", "Muffins"}
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
    misquote = misquote[:strings.LastIndex(misquote, "-")]
    misquote = misquote + "- " + fakeuser
    s.ChannelMessageSend(msg.ChannelID,misquote)

    //fmt.Println(quotes[newIndex])
    fmt.Println(misquote)
}

//Cakebombs 11/11/17
func getFake(s *discordgo.Session, msg *discordgo.MessageCreate, comp string){
    var mapping map[string][]string
    var quoteArray []string
    const MaxLen = 20;
	var count int

    mapping = make(map[string][]string, 10000)

    qte, err := db.Query("SELECT quote FROM quotes")
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
        
        quote = quote[1:strings.LastIndex(quote, "\"")]
        quotes[index] = quote
        index++
    }


    for i := 0; i<len(quotes); i++{
	    quote = quotes[i]
        quoteArray = strings.Split(quote," ")
	    for j := 0; j<len(quoteArray); j++{
            _, ok := mapping[quoteArray[j]]
			//values, ok := mapping[quoteArray[j]]
            if ok {
                if j==(len(quoteArray)-1){
                    //values = append(values, "")
					mapping[quoteArray[j]] = append(mapping[quoteArray[j]], "")
                }else{
                    //values = append(values, quoteArray[j+1])
					mapping[quoteArray[j]] = append(mapping[quoteArray[j]], quoteArray[j+1])
                }
            }else{
                if j==len(quoteArray)-1{
                    mapping[quoteArray[j]] = []string {""}
                }else{
                    mapping[quoteArray[j]] = []string {quoteArray[j+1]}
                }
            }
        }               
    }

    keys := []string {}
    for key, _ := range mapping {
        keys = append(keys, key)
    }


    next := ""
    fake := ""

    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    if index == 0{
        getQuote(s,msg," ")
        return
    }
	
    for count = 0; count < MaxLen; count++{
	    //fmt.Println(count)
        if count==0{
            values, ok := mapping[comp]
    	    if ok{
                newIndex = r1.Intn(len(values))
		        next = values[newIndex]
                fake = comp + " " + next
            }else{
                newIndex = r1.Intn(len(keys))
                get := keys[newIndex]
                values := mapping[get]
				newIndex = r1.Intn(len(values))
		        next = values[newIndex]
                fake = get + " " + next                
            }
        }else{
            values := mapping[next]
            newIndex = r1.Intn(len(values))
            next = values[newIndex]
            fake = fake + " " + next
        }
		
        if next==""{
		    if count >=6{
		        break
		    }else{
			    count=-1
			    fake = ""
		    }
		}  
		
    }
	fake = fake[:strings.LastIndex(fake, " ")]
	fake = "\"" + fake + "\"" + " - " + fakeusers[r1.Intn(usercount)]
	
	
    //fmt.Println(quotes)
    //fmt.Println(quotes[newIndex])
    s.ChannelMessageSend(msg.ChannelID,fake)
    fmt.Println(fake)
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
    CmdList["fakequote"] = getFake
}