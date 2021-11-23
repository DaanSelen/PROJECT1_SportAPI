package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var (
	dbpassword  string
	apipassword string
)

type User struct {
	ID     int    `json:"id"`
	Apikey string `json:"apikey"`
}

type Exercise struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Bodypart    string `json:"bodypart"`
	Musclegroup string `json:"musclegroup"`
	Content     string `json:"content"`
}

func initSys() {
	retrievePasswords()
	initDBConnection(dbpassword)
}

func handleError(err error) {
	if err != nil {
		log.Println("Error encountered: ", err)
	}
}

func checkApiKey(giftedKey string) bool {
	keys := getOnlyKeysFromDatabase()
	for _, key := range keys {
		if giftedKey == key {
			return true
		}
	}
	return false
}

func retrievePasswords() {
	data, err := ioutil.ReadFile("dbpassword.key")
	handleError(err)
	dbpassword = string(data)
	data, err = ioutil.ReadFile("apipassword.key")
	handleError(err)
	apipassword = string(data)
}

func checkIfInt(markedString string) bool {
	if _, err := strconv.Atoi(markedString); err == nil {
		return true
	} else {
		return false
	}
}

func getAllKeysWithID() []User {
	allKeys := retrieveAllKeysWithID()
	return allKeys
}

func createAPIKey() string {
	newKey := generateUUIDKey()
	storeKey(newKey)
	return newKey
}
func generateUUIDKey() string {
	keyWithHyphen := uuid.New()
	uuid := strings.Replace(keyWithHyphen.String(), "-", "", -1)
	return uuid
}

func deleteAPIKey(markedKey string) {
	deleteKey(markedKey) //passes the func to the database layer
}

func getAllExercises() []Exercise {
	allExercises := retrieveAllExercises()
	return allExercises
}

func getSpecificExercise(idQuery string, mgQuery string) []Exercise { //Checks status and passes on to database layer
	var exercises []Exercise
	if len(mgQuery) == 0 {
		exercises = retrieveExercise(1, idQuery)
	} else if len(idQuery) == 0 {
		exercises = retrieveExercise(2, mgQuery)
	}
	return exercises
}

func checkBAPerms(baUsername, baPassword string, baOk bool) bool {
	if !baOk {
		return false
	} else if baUsername == "admin" && baPassword == apipassword {
		return true
	}
	return false
}
