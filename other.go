package main

import (
	"database/sql"
	"fmt"
)

func getUser(db *sql.DB, path string) (UserInformation, error) {
	rows1, err := db.Query("SELECT userName, pass FROM users WHERE userName = ?", path)
	if err != nil {
		fmt.Println("1: ", err.Error())
		return UserInformation{}, err
	}
	var u UserInformation
	for rows1.Next() {
		err := rows1.Scan(&u.Username, &u.Pass)
		if err != nil {
			fmt.Println("2: ", err.Error())
			return UserInformation{}, err
		}
	}
	err = u.getInformation(db)
	if err != nil {
		return UserInformation{}, err
	}
	return u, nil
}
func (u *UserInformation) getInformation(db *sql.DB) error {
	rows, err := db.Query("SELECT userName, title,family,biography,gender FROM informations WHERE userName = ?", u.Username)
	if err != nil {
		return err
	}
	defer rows.Close()
	var ui UserInformation
	for rows.Next() {
		rows.Scan(&ui.Username, &ui.Title, &ui.Family, &ui.Biography, &ui.Gender)
	}
	u.Biography = ui.Biography
	u.Title = ui.Title
	return nil
}
