package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	csvFileName     = "data.csv"
	insertBatchSize = 50
)

type MeterReading struct {
	NMI         string
	Date        time.Time
	Consumption float64
}

func readCSVFile() ([][]string, error) {
	file, err := os.Open(csvFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	return csvReader.ReadAll()
}

func validate200Record(record []string) error {
	if len(record) != 10 {
		return errors.New("invalid number of elements in the 200 record")
	}

	// TODO: validate nmi format
	// TODO: validate interval format, should be a valid number

	return nil
}

// Sample: 200,NEM1201009,E1E2,1,E1,N1,01009,kWh,30,20050610
func parse200Record(record []string) (string, int, error) {
	if err := validate200Record(record); err != nil {
		return "", 0, err
	}

	interval, err := strconv.Atoi(record[8])
	if err != nil {
		return "", 0, err
	}

	return record[1], interval, nil
}

// Example: RecordIndicator,IntervalDate,IntervalValue1 . . . IntervalValueN,
// QualityMethod,ReasonCode,ReasonDescription,UpdateDateTime,MSATSLoadDateTime
// 300,20030501,50.1, . . . ,21.5,V,,,20030101153445,20030102023012
func validate300Record(record []string, interval int) error {
	numOfRequiredElements := ((24 * 60) / interval) + 7
	if len(record) != numOfRequiredElements {
		return errors.New("invalid number of elements in the 300 record")
	}

	// TODO: validate if all the consumption values are valid numbers
	// TODO: validate if the date is valid format

	return nil
}

/*
	 Sample: 300,20050301,0,0,0,0,0,0,0,0,0,0,0,0,0.461,0.810,0.568,1.234,1.353,1.507,1.344,1.773,0
		.848,1.271,0.895,1.327,1.013,1.793,0.988,0.985,0.876,0.555,0.760,0.938,0.566,0.512,0.9
		70,0.760,0.731,0.615,0.886,0.531,0.774,0.712,0.598,0.670,0.587,0.657,0.345,0.231,A,,,2
		0050310121004,20050310182204
*/
func parse300Record(record []string, interval int) (time.Time, []string, error) {
	if err := validate300Record(record, interval); err != nil {
		return time.Time{}, nil, err
	}

	date, err := ConvertStringToDate(record[1])
	if err != nil {
		return time.Time{}, nil, err
	}

	return date, record[2:(((24 * 60) / interval) + 2)], nil
}

func createMeterReadingList(data [][]string) ([]MeterReading, error) {
	var meterList []MeterReading
	var currNMI string
	var currIntervalInMin int // important to let the name tell about min/sec etc

	for _, line := range data {
		recordType := line[0]

		switch recordType {
		case "200":
			currNMI, currIntervalInMin, _ = parse200Record(line)

		case "300":
			// 300 can only be grouped if 200 was present before that
			if currNMI != "" {
				date, consumptionArray, _ := parse300Record(line, currIntervalInMin)
				for i, c := range consumptionArray {
					value, _ := strconv.ParseFloat(c, 64)
					meterList = append(meterList, MeterReading{
						NMI:         currNMI,
						Date:        date.Add(time.Minute * time.Duration(currIntervalInMin*(i+1))),
						Consumption: value,
					})
				}
			}

		default:
			// reset 200 related values
			currNMI = ""
			currIntervalInMin = 0
		}
	}

	return meterList, nil
}

func createBatchInsertStatements(records []MeterReading) ([]string, error) {
	var batchInsertStatements []string

	for i := 0; i < len(records); i += insertBatchSize {
		end := i + insertBatchSize
		if end > len(records) {
			end = len(records)
		}

		var values string
		for j, rec := range records[i:end] {
			if j > 0 {
				values += ","
			}
			id := uuid.New().String()
			values += fmt.Sprintf("('%s', '%s', '%s', %f)", id, rec.NMI, rec.Date.String(), rec.Consumption)
		}

		batchInsertStatements = append(batchInsertStatements, fmt.Sprintf("INSERT INTO meter_readings (id, nmi, timestamp, consumption) VALUES %s;", values))
	}

	return batchInsertStatements, nil
}

func main() {
	data, err := readCSVFile()
	if err != nil {
		log.Fatal("Error reading CSV file:", err)
	}

	meterList, err := createMeterReadingList(data)
	if err != nil {
		log.Fatal("Error creating meter reading list:", err)
	}

	insertStatements, err := createBatchInsertStatements(meterList)
	if err != nil {
		log.Fatal("Error creating batch insert statements:", err)
	}

	fmt.Printf("Insert Statements: %+v\n", insertStatements)
}
