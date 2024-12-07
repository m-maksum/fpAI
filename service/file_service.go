package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"errors"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	// return nil, nil // TODO: replace this
	if fileContent == "" {
		return nil, errors.New("CSV file is empty")
	}

	reader := csv.NewReader(strings.NewReader(fileContent))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errors.New("CSV file has no records")
	}

	headers := records[0]
	if len(headers) == 0 {
		return nil, errors.New("CSV file does not have a header row")
	}

	data := make(map[string][]string)
	for _, header := range headers {
		data[header] = []string{}
	}

	for _, record := range records[1:] {
		if len(record) != len(headers) {
			return nil, errors.New("CSV file has inconsistent number of columns")
		}
		for i, value := range record {
			data[headers[i]] = append(data[headers[i]], value)
		}
	}

	return data, nil
}
