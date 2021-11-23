package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var (
	celdrithDB *sql.DB
	err        error
)

func initDBConnection(password string) {
	celdrithDB, err = sql.Open("mysql", "celdrith:"+password+"@tcp(192.168.178.23:3306)/celdrithdb") //IF connected with OPENVPN SERVER.
	if err != nil {
		log.Println("Connection to Database Server failed: ", err)
	}
}

func retrieveAllKeysWithID() []User {
	data, err := celdrithDB.Query("SELECT * FROM user")
	handleError(err)
	defer data.Close()
	var allKeys []User
	for data.Next() {
		var key User
		err := data.Scan(&key.ID, &key.Apikey)
		handleError(err)
		allKeys = append(allKeys, key)
	}
	return allKeys
}
func getOnlyKeysFromDatabase() []string {
	data, err := celdrithDB.Query("SELECT Apikey FROM user")
	handleError(err)
	defer data.Close()
	var keys []string
	for data.Next() {
		var key string
		err := data.Scan(&key)
		handleError(err)
		keys = append(keys, key)
	}
	return keys
}

func getMaxIDCountFromUserTable() int {
	var maxID int
	data, err := celdrithDB.Query("SELECT MAX(ID) FROM user")
	handleError(err)
	defer data.Close()
	for data.Next() {
		if err := data.Scan(&maxID); err != nil {
			handleError(err)
		}
	}
	return maxID
}

func storeKey(newKey string) {
	maxID := getMaxIDCountFromUserTable()
	strMaxID := strconv.Itoa(maxID + 1)
	_, err = celdrithDB.Query("INSERT INTO user VALUES ('" + strMaxID + "', '" + newKey + "')")
	handleError(err)
}

func deleteKey(markedKey string) {
	_, err = celdrithDB.Query("DELETE FROM user WHERE ID='" + markedKey + "'")
	handleError(err)
}

func retrieveExercise(index int, input string) []Exercise {
	querySentence := "SELECT * FROM exercise WHERE "
	if index == 1 {
		querySentence += "id = " + input
	} else if index == 2 {
		querySentence += "musclegroup LIKE '%" + input + "%'"
	}
	data, err := celdrithDB.Query(querySentence)
	handleError(err)
	defer data.Close()
	var exercises []Exercise
	for data.Next() {
		var exercise Exercise
		err := data.Scan(&exercise.ID, &exercise.Title, &exercise.Bodypart, &exercise.Musclegroup, &exercise.Content)
		handleError(err)
		exercises = append(exercises, exercise)
	}
	return exercises
}

func retrieveAllExercises() []Exercise {
	data, err := celdrithDB.Query("SELECT * FROM exercise")
	handleError(err)
	defer data.Close()
	var allExercises []Exercise
	for data.Next() {
		var exercise Exercise
		err := data.Scan(&exercise.ID, &exercise.Title, &exercise.Bodypart, &exercise.Musclegroup, &exercise.Content)
		handleError(err)
		allExercises = append(allExercises, exercise)
	}
	return allExercises
}
