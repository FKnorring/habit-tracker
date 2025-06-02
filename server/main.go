package main

import (
	"log"
	"net/http"
)

/*
	Endpoints:
		GET /habits
		GET /habits/:id
		POST /habits
		PUT /habits/:id
		DELETE /habits/:id
		POST /habits/:id/tracking
		GET /habits/:id/tracking
*/

func getHabits(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getHabits stub"))
}

func getHabit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getHabit stub"))
}

func createHabit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createHabit stub"))
}

func updateHabit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateHabit stub"))
}

func deleteHabit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteHabit stub"))
}

func createTracking(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createTracking stub"))
}

func getTracking(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getTracking stub"))
}

func main() {
	http.HandleFunc("/habits", getHabits)
	http.HandleFunc("/habits/:id", getHabit)
	http.HandleFunc("/habits", createHabit)
	http.HandleFunc("/habits/:id", updateHabit)
	http.HandleFunc("/habits/:id", deleteHabit)
	http.HandleFunc("/habits/:id/tracking", createTracking)
	http.HandleFunc("/habits/:id/tracking", getTracking)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
