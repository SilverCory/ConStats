package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/SilverCory/constats/speedtest"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

// MySQL instance
type MySQL struct {
	Host string
}

const (
	tableInit = "CREATE TABLE IF NOT EXISTS {TABLE_NAME} (" +
		"`time` DATETIME NOT NULL DEFAULT NOW()," +
		"`ping` FLOAT NOT NULL DEFAULT 0," +
		"`upload` INT NOT NULL DEFAULT 0," +
		"`download` INT NOT NULL DEFAULT 0," +
		"PRIMARY KEY (`time`)," +
		"UNIQUE INDEX `time_UNIQUE` (`time` ASC));"

	insertStmt = "INSERT INTO {TABLE_NAME} (`time`, `ping`, `upload`, `download`) VALUES (?, ?, ?, ?);"
)

// Create creates an instance
func Create() *MySQL {
	return &MySQL{
		Host: "constats:alpine@/constats",
	}
}

func (m *MySQL) createConn(table string) (*sql.DB, error) {
	db, err := sql.Open("mysql", m.Host)
	if err != nil {
		return db, err
	}

	if table != "" {
		_, err := db.Exec(replaceTable(tableInit, table))
		if err != nil {
			return db, err
		}
	}

	return db, err
}
func replaceTable(input string, tableName string) string {
	return strings.Replace(input, "{TABLE_NAME}", "`"+tableName+"`", 1)
}

// Save saves the data provided or negatives if nil.
func (m *MySQL) Save(result *speedtest.TestResult, parsedTime *time.Time, table string) error {

	db, err := m.createConn(table)
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare(replaceTable(insertStmt, table))
	if err != nil {
		return err
	}

	if result == nil && parsedTime == nil {
		_, err = stmt.Exec(time.Now(), -1, -1, -1)
	} else {
		_, err = stmt.Exec(parsedTime, result.Ping, result.Upload, result.Download)
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
