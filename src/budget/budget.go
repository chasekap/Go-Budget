package budget

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
)

type WellsRow struct {
	Date string  
	Amount float64 
	Unk string 
	Note string 
	Desc string 
}

type DiscRow struct {
	TransDate string  
	PostDate string 
	Description string 
	Amount float64 
	Category string 
}

func SayHello(){
	_, err := GetStructListFromCSV("./budget/disc.csv", &DiscRow{}, true)
	if err != nil {
		fmt.Println(err)
	}
}

func GetStructListFromCSV(path string, obj any, skipHeaders bool) ([]any, error){
	
	rows := []any{}


	f, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
	defer f.Close()
	csvReader := csv.NewReader(f)
	if skipHeaders {
		_, err := csvReader.Read()
        if err != nil {
            return nil, err 
        }
	}
	for {
        row, err := csvReader.Read()
		println(row)
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err 
        }
	   
       extractedRow, err := extractRow(obj, row)
	   fmt.Println(extractedRow)
	   if err != nil {
			return nil, err
	   }

	   rows = append(rows, extractedRow)
    }
	for _, row := range rows {
		fmt.Println(row.(DiscRow).Amount)
	}
	return nil, nil 
}
func extractRow(obj any, values []string) (interface{}, error) {
	
	objCopy := reflect.New(reflect.ValueOf(obj).Elem().Type())
	s := objCopy.Elem()
	for i := 0; i < s.NumField(); i++ {
		switch s.Field(i).Type().Kind(){
		case reflect.Int: 
			conv, err := strconv.Atoi(values[i])
			if err != nil {
				return nil, err 
			}
			s.Field(i).SetInt(int64(conv))
		case reflect.String: 
			s.Field(i).SetString(values[i])
		case reflect.Bool: 
			conv, err := strconv.ParseBool(values[i])
			if err != nil {
				return nil, err 
			}
			s.Field(i).SetBool(conv)
		case reflect.Float64: 
			conv, err := strconv.ParseFloat(values[i], 64)
			if err != nil {
				return nil, err 
			}
			s.Field(i).SetFloat(conv)
		default: 
			continue 
		}
	}
	return s.Interface(), nil 
}