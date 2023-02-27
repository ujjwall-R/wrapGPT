package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/ujjwall-R/go-chat-service/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

var TOKEN string

var c *gogpt.Client
var ctx context.Context
var coll *mongo.Collection

func getReplyFromBot(prompt string) string {
	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 10,
		Prompt:    prompt,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return "Sorry, Kiara is not responding! Sorry, theres problem from ours end."
	}
	return resp.Choices[0].Text
}

func checkRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("chat service live")
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user mongodb.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
	}

	existingUser, err := mongodb.SearchUser(coll, user.Userid)
	if err == nil {
		newPrompt := fmt.Sprintf("%s Q) %s", existingUser.Prompt, user.Prompt)
		message := getReplyFromBot(newPrompt)
		newPrompt = fmt.Sprintf("%s A) %s \n", newPrompt, message)
		err := mongodb.UpdateUser(coll, existingUser, newPrompt)
		if err != nil {
			fmt.Println("Error in updation")
		}
		json.NewEncoder(w).Encode(message)
		return
	}
	message := getReplyFromBot(user.Prompt)
	user.Prompt = fmt.Sprintf("Q)%s, A)%s \n", user.Prompt, message)
	_, err = mongodb.AddUser(coll, user)
	if err != nil {
		fmt.Println("Error in Adding User")
	}
	json.NewEncoder(w).Encode(message)
	return
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	TOKEN = os.Getenv("OPENAI_TOKEN")
	c = gogpt.NewClient(TOKEN)
	ctx = context.Background()
	coll = mongodb.DB()
	r := mux.NewRouter()
	r.HandleFunc("/", checkRoute).Methods("GET")
	r.HandleFunc("/", chatHandler).Methods("POST")
	fmt.Printf("Starting server at PORT 5001\n")
	log.Fatal(http.ListenAndServe(":5001", r))
}
