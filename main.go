package main

import (
	"encoding/json"
	"net/http"
)

type Note struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var database = make(map[string]Note)

func SetJsonResp(res http.ResponseWriter, message []byte, httpCode int) {
	res.Header().Set("Content-type", "applicetion/json")
	res.WriteHeader(httpCode)
	res.Write(message)
}

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		message := []byte(`{"message":"server is up"}`)
		SetJsonResp(res, message, http.StatusOK)
	})

	http.HandleFunc("/get-notes", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`{"message": "http method salah"}`)
			SetJsonResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		var notes []Note

		for _, note := range database {
			notes = append(notes, note)
		}

		noteJson, err := json.Marshal(&notes)

		if err != nil {
			message := []byte(`{"message": "error saat melakukan parsing data"}`)
			SetJsonResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJsonResp(res, noteJson, http.StatusOK)
	})

	http.HandleFunc("/add-note", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			message := []byte(`{"message": "http method salah"}`)
			SetJsonResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		var note Note

		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&note)
		if err != nil {
			message := []byte(`{"message": "error saat melakukan parsing data"}`)
			SetJsonResp(res, message, http.StatusInternalServerError)
			return
		}

		database[note.Id] = note

		message := []byte(`{"message":"Note Baru Berhasil Ditambahkan"}`)
		SetJsonResp(res, message, http.StatusCreated)
	})

	http.HandleFunc("/get-note", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`{"message": "http method salah"}`)
			SetJsonResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message":"membutuhkan note id"}`)
			SetJsonResp(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]
		note, ok := database[id]
		if !ok {
			message := []byte(`{"message": "note tidak ditemukan"}`)
			SetJsonResp(res, message, http.StatusOK)
			return
		}

		noteJson, err := json.Marshal(&note)
		if err != nil {
			message := []byte(`{"message": "error saat melakukan parsing data"}`)
			SetJsonResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJsonResp(res, noteJson, http.StatusOK)
	})

	http.HandleFunc("/delete-notes", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			message := []byte(`{"message" : "http method salah"}`)
			SetJsonResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message":"membutuhkan note id"}`)
			SetJsonResp(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]
		note, ok := database[id]

		if !ok {
			message := []byte(`{"message": "note tidak ditemukan"}`)
			SetJsonResp(res, message, http.StatusOK)
			return
		}

		delete(database, id)

		noteJson, err := json.Marshal(&note)

		if err != nil {
			message := []byte(`{"message": "error saat melakukan parsing data"}`)
			SetJsonResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJsonResp(res, noteJson, http.StatusOK)
	})

	http.HandleFunc("/update-notes", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			message := []byte(`{"message" : "http method salah"}`)
			SetJsonResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message":"membutuhkan note id"}`)
			SetJsonResp(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]
		note, ok := database[id]

		if !ok {
			message := []byte(`{"message": "note tidak ditemukan"}`)
			SetJsonResp(res, message, http.StatusOK)
			return
		}

		var newNote Note

		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&newNote)

		if err != nil {
			message := []byte(`{"message": "error saat melakukan parsing data"}`)
			SetJsonResp(res, message, http.StatusInternalServerError)
			return
		}

		note.Title = newNote.Title
		note.Body = newNote.Body

		database[note.Id] = note

		noteJson, err := json.Marshal(&note)
		if err != nil {
			message := []byte(`{"message": "error saat melakukan parsing data"}`)
			SetJsonResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJsonResp(res, noteJson, http.StatusOK)
	})

	http.ListenAndServe(":3005", nil)
}
