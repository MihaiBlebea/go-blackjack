package persistence

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	createTable()
}

func connect() *sql.DB {
	user, password, host, database := "admin", "root", "127.0.0.1:32768", "dev_logger"
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+")/"+database+"?parseTime=true")
	checkError(err)
	return db
}

func createTable() {
	db := connect()
	defer db.Close()

	stm, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS blackjacks(
			id INT AUTO_INCREMENT PRIMARY KEY,
			dealer_hand TEXT,
			player_hand TEXT,
			dealer_score INT,
			player_score INT,
			bet INT,
			credit INT,
			result INT,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	checkError(err)

	_, err = stm.Exec()
	checkError(err)
}

func insert(bj Blackjack) int {
	db := connect()
	defer db.Close()

	stm, err := db.Prepare(`
		INSERT INTO blackjacks (
			dealer_hand,
			player_hand,
			dealer_score,
			player_score,
			bet,
			credit,
			result
		) VALUES (
			?, ?, ?, ?, ?, ?, ?
		)
	`)
	checkError(err)

	res, err := stm.Exec(bj.dealerHand, bj.playerHand, bj.dealerScore, bj.playerScore, bj.bet, bj.credit, bj.result)
	checkError(err)

	id, err := res.LastInsertId()
	checkError(err)

	return int(id)
}

func selectId() {
	db := connect()
	defer db.Close()
}
