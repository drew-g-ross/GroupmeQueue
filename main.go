package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// BOT CODE

// COMMANDS - strings.Fields to parse text
// Parse for ! as first char
// Add - make just two teammates at first
// Clear - clear team, update queue
// Show
// Help
//  After editing queue

// QUEUE
// Give each team ID, reset IDs at 5 am. Will be easier when deleting

// WEBSITE
// When updated on website, send messsage to API so it can send notifications

// ELO
// Have everyone register the name they want for ranking
var BOT_ID string = "ed7b0530736972feafd45547bb" // CHANGE TO .env

type SendMessage struct {
	bot_id string
	text   string
}

func Respond(text string) {
	words := strings.Fields(text)
	command := strings.ToLower(words[0])
	if command == "!add" {
		sendMessage("adding")
	}
}

func sendMessage(text string) {
	// Sends text to the groupchat the bot is in

	// Create a request data object.
	message := SendMessage{
		bot_id: BOT_ID,
		text:   text,
	}

	// Marshal the request data object to JSON.
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a new POST request.
	req, err := http.NewRequest("POST", "https://api.groupme.com/v3/bots/post", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the content type and accept headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Create a new HTTP client and send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read the response body.
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

// SERVER CODE
type Message struct {
	Attachments []string `json:"attachments"`
	AvatarURL   string   `json:"avatar_url"`
	CreatedAt   int      `json:"created_at"`
	GroupID     string   `json:"group_id"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	SenderID    string   `json:"sender_id"`
	SenderType  string   `json:"sender_type"`
	SourceGUID  string   `json:"source_guid"`
	System      bool     `json:"system"`
	Text        string   `json:"text"`
	UserID      string   `json:"user_id"`
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	// Check that the request method is POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Unmarshal the request body into a Message struct.
	var message Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Pass text off to bot to handle
	Respond(message.Text)

	// Write the response from the API back to the original request.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Received"}`))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	//DELETE THIS
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("messages", handleMessages)
	fmt.Println(http.ListenAndServe(":8080", nil))
	// err := http.ListenAndServe(":8080", nil)
	// if errors.Is(err, http.ErrServerClosed) {
	// 	fmt.Printf("server closed\n")
	// } else if err != nil {
	// 	fmt.Printf("error starting server: %s\n", err)
	// 	os.Exit(1)
	// }
}
