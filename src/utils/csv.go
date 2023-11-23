package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

var (
	ErrOpenCSVFile   = errors.New("open csv file failure")
	ErrCreateCSVFile = errors.New("create csv file failure")
	ErrReadCSVBytes  = errors.New("read bytes from csv file failure")
)

type csvUtil struct{}

func NewCSVUtil() *csvUtil {
	return &csvUtil{}
}

func (csvUtil) CSV2ODSlice(filePath string) ([]string, error) {
	if file, err := os.Open(filePath); err != nil {
		fmt.Println("Error while open csv file: " + err.Error())
		return nil, ErrOpenCSVFile
	} else if records, err := csv.NewReader(file).ReadAll(); err != nil {
		defer file.Close()
		fmt.Println("Error read bytes from csv file: " + err.Error())
		return nil, ErrReadCSVBytes
	} else {
		defer file.Close()
		var result []string
		for _, record := range records {
			result = append(result, record...)
		}
		return result, nil
	}
}

func (csvUtil) ODSlice2CSV(slice []string, filePath string) error {
	if file, err := os.Create(filePath); err != nil {
		fmt.Println("Error while create csv file: " + err.Error())
		return ErrCreateCSVFile
	} else {
		var result [][]string
		for _, value := range slice {
			result = append(result, []string{value})
		}
		writer := csv.NewWriter(file)
		writer.WriteAll(result)
		writer.Flush()
		return nil
	}
}
