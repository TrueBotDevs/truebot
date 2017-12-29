package main

import(
    "fmt"
    "log"
    "time"
    "math/rand"
    "strings"
    "strconv"
    "github.com/bwmarrin/discordgo"
) 


var fakeusers = [16]string{"Ed", "Cakebombs", "Kenos", "Oblivion", "TheTrooble", "Trochlis", "Church", "ZachSK", "Kirkq", "Matty", "Twinge", "Slurpee", "Sent", "z1m", "FearfulFerret", "Muffins"}
var usercount = 16
var defaultThreshold = 10

//Begin Helper Functions
func getQuoteParts(quote string)(string,string){
    var parts []string
    parts = strings.Split(quote,"\" - ")
    if(len(parts) == 1){
        parts = strings.Split(quote,"\"- ")
    }
    if(len(parts) == 1){
        parts = strings.Split(quote,"\"-")
    }
    if(len(parts) == 1){
        parts = strings.Split(quote,"\" -")
    }
    if(len(parts) == 1){
        return "error", "error"
    }
    return strings.TrimPrefix(parts[0],"\""), parts[1] 
}

func makeQuoteFromParts(quoteText string,quotee string)(string){
    return "\"" + quoteText + "\" - " + quotee
}

func convertNameFromMap(name string)(string){
    fmt.Println(strings.ToLower(name))
    if(usermap[strings.ToLower(name)] != ""){
        return usermap[strings.ToLower(name)]
    }
    return name
}

//End Helper Functions
func getQuote(s *discordgo.Session, msg *discordgo.MessageCreate, comp string){
    qte, err := db.Query("SELECT quote, quotee FROM quotes WHERE quote LIKE \"%"+comp+"%\" OR quotee LIKE \"%"+comp+"%\"")
    if err != nil {
		log.Fatal("Query error:", err)
	}
    defer qte.Close()
    
    var quoteText string
    var quotee string
    var quotes [10000]string
    var index = 0
    var newIndex = 1
    for qte.Next(){
        err = qte.Scan(&quoteText, &quotee)
        if err != nil {
            log.Fatal("Parse error:", err)
        }
        quotes[index] = makeQuoteFromParts(quoteText,quotee)
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
    qte, err := db.Query("SELECT quote FROM quotes WHERE quote LIKE \"%"+comp+"%\" OR quotee LIKE \"%"+comp+"%\"")
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

    userchoice = r1.Intn(usercount)
    fakeuser = fakeusers[userchoice]
    misquote = makeQuoteFromParts(quotes[newIndex], fakeuser) 
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
        
        //quote = quote[1:strings.LastIndex(quote, "\"")]
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
    quoteText, quotee := getQuoteParts(quote)
    quotee = convertNameFromMap(quotee)
    if quoteText == " "{
        s.ChannelMessageSend(msg.ChannelID,"Usage: !addquote <quote>")
    }else if quoteText == "<quote>"{
        s.ChannelMessageSend(msg.ChannelID,"Very funny Church")
    }else if quoteText != "error"{
        if strings.Contains(quote, "<@"){
            s.ChannelMessageSend(msg.ChannelID, "Fuck you, don't @ people in quotes")
        }else{
            newItem := "INSERT INTO quotes (quote,quotee,submitter,date) values (?,?,?,?)"
            stmt, err := db.Prepare(newItem)
            if err != nil { panic(err) }
            defer stmt.Close()

            _, err2 := stmt.Exec(quoteText,quotee,msg.Author.Username,time.Now().Unix())
            if err2 != nil { panic(err2) }
            s.ChannelMessageSend(msg.ChannelID, "Added your quote to the database:```" + makeQuoteFromParts(quoteText,quotee) + "```")
        }
    }else{
        s.ChannelMessageSend(msg.ChannelID, "Quotes should be in the format:```\"The thing that was said\" - Username```")
    }
}

func quoteLeaderboard(s *discordgo.Session, msg *discordgo.MessageCreate, quote string){
	var threshold int
	quote = strings.TrimSpace(quote)
	if(len(quote) <= 1){
		threshold = defaultThreshold
	}else{
	    thresh, err := strconv.Atoi(quote)
		threshold = thresh
		if err != nil {
			threshold = defaultThreshold
			s.ChannelMessageSend(msg.ChannelID,"Invalid Threshold")
		}
	}
	fmt.Println(strconv.Itoa(threshold))
    qte, err := db.Query("SELECT DISTINCT quotee, COUNT(quotee) AS CountOf FROM quotes GROUP BY quotee HAVING CountOf >= "+strconv.Itoa(threshold)+" ORDER BY CountOf DESC, quotee ASC ")
    if err != nil {
		log.Fatal("Query error:", err)
	}
    defer qte.Close()
    
    var quotee string
    var quoteCount int
    var quotees [10000]string
    var quoteCounts [10000]int
    var index = 0
    
    var outputTable = "BGC Quote Leaderboard\n```"
    
    for qte.Next(){
        err = qte.Scan(&quotee, &quoteCount)
        if err != nil {
            log.Fatal("Parse error:", err)
        }
        quotees[index] = quotee
        quoteCounts[index] = quoteCount
        index++
    }
    for i:=0; i<index; i++{
        outputTable += strconv.Itoa(quoteCounts[i]) + " - " + quotees[i] + "\n"
    }
    outputTable += "```"
    s.ChannelMessageSend(msg.ChannelID,outputTable)
}

func removeQuote(s *discordgo.Session, msg *discordgo.MessageCreate, quote string){

}

func init() {
    CmdList["misquote"] = misQuote
    CmdList["quote"] = getQuote
    CmdList["addquote"] = addQuote
    CmdList["quoteLeaderboard"] = quoteLeaderboard
	AliasList["ql"] = quoteLeaderboard
    CmdList["fakequote"] = getFake
}