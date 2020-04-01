package utils

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/feed", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		challenges, ok := r.URL.Query()["hub.challenge"]

		if !ok || len(challenges[0]) < 1 {
			log.Println("Url Param 'challenge' is missing")
			return
		}
		challenge := challenges[0]

		log.Println("Url Param 'challenge' is: " + string(challenge))
		_, err := w.Write([]byte(challenge))
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == http.MethodPost {

	}
}
