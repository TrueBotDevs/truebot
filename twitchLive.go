package main

import(
    "fmt"
    "log"
    "github.com/bwmarrin/discordgo"
    "net/http"
	"io/ioutil"
	"encoding/json"
	"time"
) 

type Stream struct {
	StreamData []Data `json:"data"`	
	StreamPagination Pagination `json:"pagination"`
}

type Data struct{
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	GameID       string    `json:"game_id"`
	CommunityIds []string  `json:"community_ids"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

type Pagination struct{
	Cursor string `json:"cursor"`
}



// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
//func twitchLive(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
func twitchLive(arg string){
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login=" + arg, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Client-Id", "gddz9zfpx3zhgrexzo6rbjjntnz5y3r")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}

	body, err := ioutil.ReadAll(resp.Body)
	
	defer resp.Body.Close()
		
	if err != nil {
		log.Fatal("TwitchAPI error:", err)
	}
	
	stream, err := getStreams([]byte(body))

	if len(stream.StreamData) > 0{
		current := time.Now().Unix()
		start:= stream.StreamData[0].StartedAt.Unix()
		if stream.StreamData[0].Type == "live" && current-start < 60{
		//fmt.Println(stream.StreamData[0].Type)
			dgSession.ChannelMessageSend("362408790051651597","https://www.twitch.tv/" + arg+" " + stream.StreamData[0].Title)
		}
	}
	//fmt.Println(stream.Data)
}

func getStreams(body []byte) (*Stream, error) {
    var s = new(Stream)
    err := json.Unmarshal(body, &s)
    if(err != nil){
        fmt.Println("whoops:", err)
    }
    return s, err
}

func addStream(s *discordgo.Session, msg *discordgo.MessageCreate, stream string){
    newItem := "INSERT INTO streams (TwitchPage,DiscorduID) values (?,?)"
    stmt, err := db.Prepare(newItem)
    if err != nil { panic(err) }
    defer stmt.Close()

    _, err2 := stmt.Exec(stream,msg.Author.ID)
    if err2 != nil { panic(err2) }
    s.ChannelMessageSend(msg.ChannelID, "Added your stream to the database:```https:\\\\www.twitch.tv\\" + stream + "``` " )
    
}

func checkDB(){
	for true{
        if hasSession{
			qte, err := db.Query("SELECT TwitchPage FROM streams")
			if err != nil {
				log.Fatal("Query error:", err)
			}
			defer qte.Close()
			
			var stream string
			for qte.Next(){
				err = qte.Scan(&stream)
				if err != nil {
					log.Fatal("Parse error:", err)
				}				
			}
			twitchLive(stream)
			time.Sleep(60000 * time.Millisecond)
		}
	}
	
}

func init() {
	//CmdList["twitchtest"] = twitchLive
	CmdList["addStream"] = addStream
	go checkDB()
}
//TODO
/*

*/