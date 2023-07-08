package converters

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func ConvertExcelToJSON(filepath string, hasHeaders bool) []byte {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var objects []map[string]string
	var headers []string

	rows, err := f.GetRows("Sheet1")

	if hasHeaders {
		headers = rows[0]
	}

	if err != nil {
		return nil
	}
	obj := make(map[string]string)
	for _, row := range rows {
		for i, colCell := range row {
			if hasHeaders {
				headerCell := headers[i]
				obj[headerCell] = colCell
			} else {
				obj[strconv.Itoa(i+1)] = colCell
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
