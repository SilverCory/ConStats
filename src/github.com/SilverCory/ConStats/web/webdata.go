package web

import (
	"fmt"
	"time"

	"github.com/SilverCory/ConStats/sql"
)

type Data struct {
	Type       string      `json:"type"`
	DataPoints []DataPoint `json:"dataPoints"`
	XValueType string      `json:"xValueType"`
	Name       string      `json:"name"`
}

type DataPoint struct {
	X interface{} `json:"x"`
	Y interface{} `json:"y"`
}

type rawTime []byte

func (t rawTime) Time() (interface{}, error) {
	timeOut, err := time.Parse("2006-01-02 15:04:05", string(t))
	if err != nil {
		return -1, err
	}

	return timeOut.Unix() * 1000, nil

}

func GenerateData(storage *sql.MySQL) (*[]Data, error) {

	PingData := Data{
		Type:       "line",
		XValueType: "dateTime",
		Name:       "Ping",
		DataPoints: make([]DataPoint, 0),
	}

	UploadData := Data{
		Type:       "line",
		XValueType: "dateTime",
		Name:       "Up",
		DataPoints: make([]DataPoint, 0),
	}

	DownloadData := Data{
		Type:       "line",
		XValueType: "dateTime",
		Name:       "Down",
		DataPoints: make([]DataPoint, 0),
	}

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

		PingData.DataPoints = append(PingData.DataPoints, DataPoint{
			X: unixTime,
			Y: ping,
		})
		UploadData.DataPoints = append(UploadData.DataPoints, DataPoint{
			X: unixTime,
			Y: upload,
		})

		DownloadData.DataPoints = append(DownloadData.DataPoints, DataPoint{
			X: unixTime,
			Y: download,
		})

	}

	return &[]Data{PingData, UploadData, DownloadData}, nil

}
