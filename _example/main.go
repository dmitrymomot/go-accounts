package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	accounts "github.com/dmitrymomot/go-accounts"
	"github.com/dmitrymomot/random"
	_ "github.com/mattn/go-sqlite3" // or another driver
)

var (
	uid = "4086bf60-112b-11ea-8323-075003f81360"
)

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`
	CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		disabled INTEGER,
		created_at INTEGER NOT NULL,
		updated_at INTEGER
	);

	CREATE TABLE IF NOT EXISTS members (
		id TEXT PRIMARY KEY,
		account_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at INTEGER NOT NULL
	);
	`)

	ai := accounts.NewInteractor(
		accounts.NewAccountRepository(db, "sqlite3", "accounts"),
		accounts.NewMemberRepository(db, "sqlite3", "members"),
	)

	if err := ai.Insert(&accounts.Account{Name: "The App " + random.String(5)}, uid); err != nil {
		panic(err)
	}
	if err := ai.Insert(&accounts.Account{Name: "The App " + random.String(5)}, uid); err != nil {
		panic(err)
	}
	if err := ai.Insert(&accounts.Account{Name: "The App " + random.String(5)}, uid); err != nil {
		panic(err)
	}

	al, err := ai.GetAccountsListWithRoleByUserID(uid,
		accounts.OrderBy(accounts.OrderByCreatedAtDesc),
		accounts.Limit(3),
		accounts.Role(accounts.RoleOwner),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\naccounts in list: %d \n", len(al))

	for _, a := range al {
		fmt.Println(fmt.Printf("role: %s; account: %s, created_at: %s", a.Role, a.Name, time.Unix(a.CreatedAt, 0).String()))
	}
}
