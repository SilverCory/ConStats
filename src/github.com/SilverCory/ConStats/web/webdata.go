package web

import (
	"fmt"
	"time"

	"strings"

	"strconv"

	"github.com/SilverCory/ConStats/sql"
)

// Data the data to show.
type Data struct {
	Type         string      `json:"type"`
	DataPoints   []DataPoint `json:"dataPoints"`
	XValueType   string      `json:"xValueType"`
	Name         string      `json:"name"`
	Unit         string      `json:"unit"`
	ShowInLegend bool        `json:"showInLegend"`
}

// DataPoint the point of data to display.
type DataPoint struct {
	TimePoint int
	Ping      float32
	Upload    float32
	Download  float32
}

type rawTime []byte

func (t rawTime) Time() (string, error) {
	timeOut, err := time.Parse("2006-01-02 15:04:05", string(t))
	if err != nil {
		return "", err
	}

	stringDate := timeOut.Format("Date(2006,01,02,15,04,05)")

	// Prepare for fucking gross shit.
	dateParts := strings.Split(stringDate, ",")

	monthInt, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return "", err
	}

	dateParts[1] = strconv.Itoa(monthInt - 1)

	return strings.Join(dateParts, ","), err

}

// GenerateData generates the data statistics.
func GenerateData(storage *sql.MySQL) ([]interface{}, error) {

	interaceArray := make([]interface{}, 0)

	rows, err := storage.Load()
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var timePoint rawTime
		var ping float32
		var upload float32
		var download float32

		err := rows.Scan(&timePoint, &ping, &upload, &download)
		if err != nil {
			fmt.Println("Error for a row!", err)
			continue
		}

		unixTime, err := timePoint.Time()
		if err != nil {
			fmt.Println("Error for a row!", err)
			continue
		}

		if upload > 0 {
			upload = upload / 1000000
		}

		if download > 0 {
			download = download / 1000000
		}

		interaceArray = append(interaceArray, []interface{}{unixTime, ping, upload, download})

	}

	return interaceArray, nil

}
