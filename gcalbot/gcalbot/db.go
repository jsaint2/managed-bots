package gcalbot

import (
	"database/sql"

	"github.com/keybase/managed-bots/base"
)

type DB struct {
	*base.GoogleOAuthDB
}

func NewDB(db *sql.DB) *DB {
	return &DB{
		GoogleOAuthDB: base.NewGoogleOAuthDB(db),
	}
}

func (d *DB) GetAccountsForUser(username string) (accounts []string, err error) {
	rows, err := d.DB.Query(`SELECT nickname
		FROM accounts
		WHERE username = ?
		ORDER BY nickname`, username)
	if err == sql.ErrNoRows {
		return accounts, nil
	} else if err != nil {
		return nil, err
	}
	for rows.Next() {
		var account string
		err = rows.Scan(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (d *DB) ExistsAccountForUser(username string, nickname string) (exists bool, err error) {
	row := d.DB.QueryRow(`SELECT EXISTS(
		SELECT * FROM accounts WHERE username = ? AND nickname = ?)`,
		username, nickname)
	err = row.Scan(&exists)
	return exists, err
}

func (d *DB) InsertAccountForUser(username string, nickname string) error {
	err := d.RunTxn(func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO accounts
		(username, nickname)
		VALUES (?, ?)
	`, username, nickname)
		return err
	})
	return err
}

func (d *DB) DeleteAccountForUser(username string, nickname string) error {
	identifier := GetAccountIdentifier(username, nickname)
	err := d.RunTxn(func(tx *sql.Tx) error {
		_, err := tx.Exec(`DELETE FROM oauth
	WHERE identifier = ?`, identifier)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`DELETE FROM accounts
	WHERE username = ? and nickname = ?`, username, nickname)
		return err
	})
	return err
}
