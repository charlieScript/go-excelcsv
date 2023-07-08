package converters

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

func ConvertCSVToJSON(filepath string, hasHeaders bool) []byte {
	// Open the CSV file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	var header []string

	if hasHeaders {
		header, _ = reader.Read()
	}

	if err != nil {
		log.Fatal(err)
	}
	var objects []map[string]string
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		obj := make(map[string]string)
		for i, value := range row {
			if hasHeaders {
				obj[header[i]] = value
			} else {
				obj[strconv.Itoa(i+1)] = value
			}
		}

		objects = append(objects, obj)
	}
	jsonData, err := json.Marshal(objects)
	if err != nil {
		log.Fatal(err)
	}

	return jsonData
}
