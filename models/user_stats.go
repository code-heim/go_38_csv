package models

import (
	"log"
	"strconv"
	"time"
)

type UserStats struct {
	UserID     int
	UserName   string
	LoginCount int
	LastLogin  time.Time
	Active     bool
}

func GetUserStats() *[]UserStats {
	// Query the data
	rows, err := DB.Query("SELECT UserID, UserName, LoginCount, LastLogin, Active FROM user_stats")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Slice to hold the user stats
	var users []UserStats

	// Iterate over the rows
	for rows.Next() {
		var user UserStats
		var lastLogin string

		// Scan the row into the struct fields
		err := rows.Scan(&user.UserID, &user.UserName, &user.LoginCount, &lastLogin, &user.Active)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the LastLogin time
		user.LastLogin, err = time.Parse(time.RFC3339, lastLogin)
		if err != nil {
			log.Fatal(err)
		}

		// Append to the slice
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &users
}

func AddRows(rows *[][]string) error {
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO user_stats (UserID, UserName, LoginCount, LastLogin, Active) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	for _, record := range (*rows)[1:] {
		loginCount, _ := strconv.Atoi(record[2])
		active, _ := strconv.ParseBool(record[4])

		_, err := stmt.Exec(record[0], record[1], loginCount, record[3], active)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting data: %v", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}
	return nil
}
