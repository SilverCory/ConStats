package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/SilverCory/constats/speedtest"
	_ "github.com/go-sql-driver/mysql"
)

// MySQL instance
type MySQL struct {
	Host string
}

const (
	tableInit = "CREATE TABLE IF NOT EXISTS ? (\n" +
		"`time` DATETIME NOT NULL DEFAULT NOW()," +
		"`ping` FLOAT NOT NULL DEFAULT 0," +
		"`upload` INT NOT NULL DEFAULT 0," +
		"`download` INT NOT NULL DEFAULT 0," +
		"PRIMARY KEY (`time`)," +
		"UNIQUE INDEX `time_UNIQUE` (`time` ASC));"

	insertStmt = "INSERT INTO ? (`time`, `ping`, `upload`, `download`) VALUES (?, ?, ?, ?);"
)

// Create creates an instance
func Create() *MySQL {
	return &MySQL{
		Host: "constats:alpine@/constats",
	}
}

func (m *MySQL) createConn(table string) (*sql.DB, error) {
	db, err := sql.Open("sql", m.Host)
	if err != nil {
		return db, err
	}

	if table != "" {
		stmt, err := db.Prepare(tableInit)
		if err != nil {
			return db, err
		}

		stmt.Exec(table)
	}

	return db, err
}

// Save saves the data provided or negatives if nil.
func (m *MySQL) Save(result *speedtest.TestResult, parsedTime *time.Time, table string) error {

	db, err := m.createConn(table)
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare(insertStmt)
	if err != nil {
		return err
	}

	if result == nil && parsedTime == nil {
		stmt.Exec(time.Now(), -1, -1, -1)
	} else {
		_, err = stmt.Exec(table, parsedTime, result.Ping, result.Upload, result.Download)
	}

	return err

}

func (m *MySQL) Load(table string) (*sql.Rows, error) {

	db, err := m.createConn(table)
	if err != nil {
		return nil, err
	}

	stats, err := db.Query("SELECT * FROM `" + table + "` LIMIT 3000")
	if err != nil {
		return stats, err
	}

	return stats, nil

}

func (m *MySQL) FindTables() ([]string, error) {

	db, err := m.createConn("")
	if err != nil {
		return nil, err
	}

	returnStrings := make([]string, 0)

	rows, err := db.Query("SHOW TABLES LIKE 'constats\\_%'")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			fmt.Println("Error for a table fetch row!", err)
			continue
		}

		returnStrings = append(returnStrings, tableName)
	}

	return returnStrings, nil

}
