package services
import (
	"bytes"
	"encoding/json"
	"net/http"
)
type ToxicityResponse struct {
	Score float64 `json:"score"`	
}
func AnalyzeToxicity(content string) (float64, bool, error) {
	apiURL := "http://localhost:5000/analyze"

	requestBody := map[string]string{
			"text": content,
	}
	jsonData, _ := json.Marshal(requestBody)
	resp, err:= http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, false, err
	}
	defer resp.Body.Close()
	var result ToxicityResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return 0, false, err
    }
	isFlagged := result.Score >= 0.5

    return result.Score, isFlagged, nil

}