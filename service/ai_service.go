package service

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"
	"bytes"
	"encoding/json"
	"errors"
	// "io/ioutil"
	// "log"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	// return "", nil
	// // return nil, nil // TODO: replace this
    if len(table) == 0 {
        return "", errors.New("data table is empty")
    }

    payload := map[string]interface{}{
        "table": table,
        "query": query,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", "https://api.huggingface.co/models/tapas-base-finetuned-wtq", bytes.NewBuffer(body))
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := s.Client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", errors.New("failed to get a valid response from AI model")
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    cells, ok := result["cells"].([]interface{})
    if !ok || len(cells) == 0 {
        return "", errors.New("invalid response format from AI model")
    }

    answer, ok := cells[0].(string)
    if !ok {
        return "", errors.New("invalid response format from AI model")
    }

    return answer, nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	// TODO: answer here
	// return model.ChatResponse{}, nil
	payload := map[string]interface{}{
        "context": context,
        "query":   query,
        "token":   token,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return model.ChatResponse{}, err
    }

    req, err := http.NewRequest("POST", "https://api.huggingface.co/models/microsoft/Phi-3.5-mini-instruct", bytes.NewBuffer(body))
    if err != nil {
        return model.ChatResponse{}, err
    }
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := s.Client.Do(req)
    if err != nil {
        return model.ChatResponse{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return model.ChatResponse{}, errors.New("failed to get a valid response from AI model")
    }

    var chatResponses []model.ChatResponse 
    if err := json.NewDecoder(resp.Body).Decode(&chatResponses); err != nil {
        return model.ChatResponse{}, err
    }

    if len(chatResponses) == 0 {
        return model.ChatResponse{}, errors.New("no responses from AI model")
    }

    return chatResponses[0], nil 
}
