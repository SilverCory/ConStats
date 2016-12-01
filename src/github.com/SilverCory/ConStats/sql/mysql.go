package sql

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL Driver

	"github.com/SilverCory/ConStats/speedtest"
)

// MySQL instance
type MySQL struct {
	Host string
}

const (
	tableInit = "CREATE TABLE IF NOT EXISTS `stats` (\n" +
		"`time` DATETIME NOT NULL DEFAULT NOW()," +
		"`ping` FLOAT NOT NULL DEFAULT 0," +
		"`upload` INT NOT NULL DEFAULT 0," +
		"`download` INT NOT NULL DEFAULT 0," +
		"PRIMARY KEY (`time`)," +
		"UNIQUE INDEX `time_UNIQUE` (`time` ASC));"

	insertStmt = "INSERT INTO `stats` (`time`, `ping`, `upload`, `download`) VALUES (?, ?, ?, ?);"
)

// Create creates an instance
func Create() *MySQL {
	return &MySQL{
		Host: "constats:alpine@/constats",
	}
}

func (m *MySQL) createConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", m.Host)
	if err != nil {
		return db, err
	}

	db.Exec(tableInit)
	return db, err
}

// Save saves the data provided or negatives if nil.
func (m *MySQL) Save(result *speedtest.TestResult, parsedTime *time.Time) error {

	db, err := m.createConn()
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
		_, err = stmt.Exec(parsedTime, result.Ping, result.Upload, result.Download)
	}

	return err

}

func (m *MySQL) Load() (*sql.Rows, error) {

	db, err := m.createConn()
	if err != nil {
		return nil, err
	}

	stats, err := db.Query("SELECT * FROM `stats` LIMIT 3000")
	if err != nil {
		return stats, err
	}

	return stats, nil

}
