package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("API Initialising...")
	initSys()

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/keys", handleShowAllKeys).Methods("GET")
	myRouter.HandleFunc("/key", handleGenerateKey).Methods("GET")
	myRouter.HandleFunc("/key", handleDeleteKey).Methods("DELETE")
	myRouter.HandleFunc("/exercises", handleShowAllExercises).Methods("GET")
	myRouter.HandleFunc("/exercise", handleGetExercise).Methods("GET")

	http.ListenAndServe(":61907", myRouter)
}

func handleShowAllKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	baUsername, baPassword, baOk := r.BasicAuth()
	permitted := checkBAPerms(baUsername, baPassword, baOk)
	if permitted {
		allKeys := getAllKeysWithID()
		json.NewEncoder(w).Encode(allKeys)
	} else {
		w.WriteHeader(401)
	}
}

func handleGenerateKey(w http.ResponseWriter, r *http.Request) { //admin func
	w.Header().Set("Content-Type", "application/json")

	baUsername, baPassword, baOk := r.BasicAuth()
	permitted := checkBAPerms(baUsername, baPassword, baOk)
	if !permitted {
		w.WriteHeader(401)
	}
	if permitted {
		newKey := createAPIKey() //Generate UUID and return it into newKey
		json.NewEncoder(w).Encode("New key generated: " + newKey)
	}
}

func handleDeleteKey(w http.ResponseWriter, r *http.Request) { //Admin func
	w.Header().Set("Content-Type", "application/json")

	idQuery, ok := r.URL.Query()["id"]
	if !ok || len(idQuery[0]) < 1 || idQuery[0] == "0" {
		w.WriteHeader(400)
	} else {
		baUsername, baPassword, baOk := r.BasicAuth()
		permitted := checkBAPerms(baUsername, baPassword, baOk)
		correctType := checkIfInt(idQuery[0])
		if !permitted {
			w.WriteHeader(401)
		} else if !correctType {
			w.WriteHeader(400)
		}
		if permitted && correctType {
			deleteAPIKey(idQuery[0])
			json.NewEncoder(w).Encode("Key with ID: " + idQuery[0] + " has been requested to be deleted.")
		}
	}
}

func handleGetExercise(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	keyQuery, ok := r.URL.Query()["key"]
	if !ok || len(keyQuery[0]) < 1 {
		w.WriteHeader(400)
	} else {
		if checkApiKey(keyQuery[0]) {
			var idPresent bool
			var mgPresent bool
			var exercises []Exercise
			idQuery, ok := r.URL.Query()["id"]
			if !ok || len(idQuery[0]) < 1 {
				idPresent = false
			} else {
				idPresent = true
			}
			mgQuery, ok := r.URL.Query()["mg"]
			if !ok || len(mgQuery[0]) < 1 {
				mgPresent = false
			} else {
				mgPresent = true
			}
			if idPresent && !mgPresent {
				exercises = getSpecificExercise(idQuery[0], "")
			} else if !idPresent && mgPresent {
				exercises = getSpecificExercise("", mgQuery[0])
			} else {
				w.WriteHeader(400)
			}
			json.NewEncoder(w).Encode(exercises)
		} else {
			w.WriteHeader(401)
		}
	}
}

func handleShowAllExercises(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	keyQuery, ok := r.URL.Query()["key"]
	if !ok || len(keyQuery[0]) < 1 {
		w.WriteHeader(400)
	} else {
		if checkApiKey(keyQuery[0]) {
			allExercises := getAllExercises()
			json.NewEncoder(w).Encode(allExercises)
		} else {
			w.WriteHeader(401)
		}
	}
}
