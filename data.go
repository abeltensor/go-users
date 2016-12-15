package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db, _ = sql.Open("sqlite3", "user.sqlite3")

func saveData(u *User) error {
	db.Exec("create table if not exists users (uuid text not null, firstname text, lastname test, username text, email text, password text, PRIMARY KEY(uuid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (uuid, firstname, lastname, username, email, password) values (?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.Uuid, u.FName, u.LName, u.UserName, u.Email, u.Password)
	tx.Commit()
	return err
}

func loadUser(uname, pw string) (*User, error) {
	u := &User{}
	q, err := db.Query("select username, password from users where username = '" + uname + "' Desc limit 1")
	if err != nil {
		return nil, err
	}
	for q.Next() {
		q.Scan(&u.UserName, &u.Password)
	}
	return u, nil
}
