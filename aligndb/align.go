package aligndb

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	dbh *sql.DB
)

func init() {
	var err error
	dbh, err = sql.Open("sqlite3", "file:./database.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}

	table := `CREATE TABLE IF NOT EXISTS align (
	appid INTEGER PRIMARY KEY, 
	align VARCHAR(20)
	)`

	query, err := dbh.Prepare(table)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	_, _ = query.Exec()
}

func SetAlign(id uint64, a string) {
	query, err := dbh.Prepare(`INSERT OR REPLACE INTO align (appid, align) VALUES (?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	_, _ = query.Exec(id, a)

}

func GetAlign(id uint64) string {
	var result string

	if id == 0 {
		return "absolute-center"
	}

	row := dbh.QueryRow(`SELECT align FROM align WHERE appid = ?`, id)
	err := row.Scan(&result)
	if err != nil {
		return "left"
	}
	return result
}

func Close() {
	_ = dbh.Close()
}
