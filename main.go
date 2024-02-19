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

type meterReadingRecord struct {
	nmi         string
	date        string
	consumption float64
}

func readCSVFile() ([][]string, error) {
	// open file
	f, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1 // since all lines in csv do not have equal elements
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func validate200Record(record []string) (bool, error) {
	if len(record) != 10 {
		return false, errors.New("not enough elements in the 200 record")
	}

	// validate nmi format
	// validate interval format, should be a valid number

	return true, nil
}

// Sample: 200,NEM1201009,E1E2,1,E1,N1,01009,kWh,30,20050610
func parse200Record(record []string) (string, int, error) {
	ok, err := validate200Record(record)
	if !ok {
		return "", 0, err
	}

	value, err := strconv.Atoi(record[8])
	if err != nil {
		return "", 0, err
	}

	return record[1], value, nil
}

func validate300Record(record []string, interval int) (bool, error) {
	// Example: RecordIndicator,IntervalDate,IntervalValue1 . . . IntervalValueN,
	// QualityMethod,ReasonCode,ReasonDescription,UpdateDateTime,MSATSLoadDateTime
	// 300,20030501,50.1, . . . ,21.5,V,,,20030101153445,20030102023012

	numOfRequiredElements := ((24 * 60) / interval) + 7
	fmt.Printf("length = %d, and required = %d", len(record), numOfRequiredElements)

	if len(record) != numOfRequiredElements {
		fmt.Printf("record = %+v", record)
		return false, errors.New("not enough elements in the 300 record")
	}

	return true, nil
}

/*
	 Sample: 300,20050301,0,0,0,0,0,0,0,0,0,0,0,0,0.461,0.810,0.568,1.234,1.353,1.507,1.344,1.773,0
		.848,1.271,0.895,1.327,1.013,1.793,0.988,0.985,0.876,0.555,0.760,0.938,0.566,0.512,0.9
		70,0.760,0.731,0.615,0.886,0.531,0.774,0.712,0.598,0.670,0.587,0.657,0.345,0.231,A,,,2
		0050310121004,20050310182204
*/
func parse300Record(record []string, interval int) (time.Time, []string, error) {
	ok, err := validate300Record(record, interval)
	if !ok {
		return time.Time{}, nil, err
	}

	date, err := ConvertStringToDate(record[1])
	if err != nil {
		return time.Time{}, nil, err
	}

	return date, record[2:(((24 * 60) / interval) + 2)], nil
}

func createMeterReadingList(data [][]string) ([]meterReadingRecord, error) {
	var meterReadingList []meterReadingRecord
	var currNmi string
	var currIntervalInMin int // important to let the name tell about min/sec etc
	var err error

	for _, line := range data {
		recordType := line[0] // 200 | 300

		switch recordType {
		case "200":
			currNmi, currIntervalInMin, err = parse200Record(line)
			if err != nil {
				return nil, err
			}

		case "300":
			// 300 can only be grouped if 200 was present before that
			if currNmi != "" {
				date, consumptionArray, err := parse300Record(line, currIntervalInMin)
				if err != nil {
					return nil, err
				}

				for i, c := range consumptionArray {
					// convert consumption string to float64
					value, err := strconv.ParseFloat(c, 64)
					if err != nil {
						return nil, err
					}

					meterReadingList = append(meterReadingList, meterReadingRecord{
						nmi:         currNmi,
						date:        date.Add(time.Minute * time.Duration(currIntervalInMin*(i+1))).String(),
						consumption: value,
					})
				}
			}

		default:
			// reset 200 related values
			currNmi = ""
			currIntervalInMin = 0
			err = nil
		}
	}

	return meterReadingList, nil
}

func createBatchInsertStatements(records []meterReadingRecord) ([]string, error) {
	const INSERT_BATCH_SIZE = 50
	var batchInsertStatements []string
	var insertStatement string

	totalRecords := len(records)
	for i := 0; i < totalRecords; i += INSERT_BATCH_SIZE {
		end := i + INSERT_BATCH_SIZE
		if end > totalRecords {
			end = totalRecords
		}

		insertStatement = "INSERT INTO meter_readings (id, nmi, timestamp, consumption) VALUES "
		for j, rec := range records[i:end] {
			id := uuid.New().String()
			if j > 0 {
				insertStatement += ","
			}
			insertStatement += fmt.Sprintf("('%s','%s', '%s', %f)", id, rec.nmi, rec.date, rec.consumption)
		}
		insertStatement += ";"
		batchInsertStatements = append(batchInsertStatements, insertStatement)
	}

	return batchInsertStatements, nil
}

func main() {
	println("starting program")

	data, err := readCSVFile()
	if err != nil {
		fmt.Println("Error1")
		log.Fatal(err)
	}

	meterList, err := createMeterReadingList(data)
	if err != nil {
		fmt.Println("Error2")
		log.Fatal(err)
	}

	insertStatement, err := createBatchInsertStatements(meterList)
	if err != nil {
		fmt.Println("Error3")
		log.Fatal(err)
	}

	fmt.Printf("insertStatement = %+v", insertStatement[0])
}
