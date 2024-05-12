package servers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"hermes/auth"
	"hermes/database"

	"github.com/tidwall/gjson"
)

const SERVER_PORT = "3333"

func postToTopic(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body")
	}

	parsed_body := gjson.Parse(string(body))
	topic := parsed_body.Get("topic")

	if topic.Str != "" {
		database.InsertObject(topic.Str, parsed_body)
	}
}

func handlePublishRequest(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if auth.IsAuthenticated(r) {
		postToTopic(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func handlePublisherRegister(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	register_key := r.Header.Get("REGISTERKEY")
	if register_key == "my_register_key" {
		uuid := auth.RegisterAsPublisher()
		io.WriteString(w, uuid)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}


func InitHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/publish", handlePublishRequest)
	mux.HandleFunc("/register", handlePublisherRegister)

	err := http.ListenAndServe(":"+SERVER_PORT, mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("[HTTP-SERVER] - %v - Server running on port %v", time.Now(), SERVER_PORT)
	}

}
