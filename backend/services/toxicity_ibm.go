package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// IBMRequest matches the IBM MAX model API format
type IBMRequest struct {
    Text []string `json:"text"`
}

// IBMResponse matches the actual IBM MAX model response
type IBMResponse struct {
    Status  string      `json:"status"`
    Results []IBMResult `json:"results"`
}

// IBMResult contains toxicity scores for one text
type IBMResult struct {
    OriginalText string `json:"original_text"`
    Predictions  struct {
        Toxic        float64 `json:"toxic"`
        SevereToxic  float64 `json:"severe_toxic"`
        Obscene      float64 `json:"obscene"`
        Threat       float64 `json:"threat"`
        Insult       float64 `json:"insult"`
        IdentityHate float64 `json:"identity_hate"`
    } `json:"predictions"`
}

// IBMAnalyzer handles communication with IBM MAX model
type IBMAnalyzer struct {
    modelURL string
    client   *http.Client
}

// NewIBMAnalyzer creates a new IBM model client
func NewIBMAnalyzer() *IBMAnalyzer {
    return &IBMAnalyzer{
        modelURL: "http://localhost:5000/model/predict",
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// Analyze sends text to IBM model and returns toxicity scores
func (i *IBMAnalyzer) Analyze(content string) (*ToxicityResult, error) {
    // Prepare request
    reqBody := IBMRequest{
        Text: []string{content},
    }
    
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %v", err)
    }

    log.Printf("📤 Sending to IBM: %s", content)

    // Make request to IBM model
    resp, err := i.client.Post(i.modelURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to IBM model: %v", err)
    }
    defer resp.Body.Close()

    // Check response status
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("IBM model returned error status: %d", resp.StatusCode)
    }

    // Parse response
    var ibmResp IBMResponse
    err = json.NewDecoder(resp.Body).Decode(&ibmResp)
    if err != nil {
        return nil, fmt.Errorf("failed to parse IBM response: %v", err)
    }

    log.Printf("📥 Received from IBM: %+v", ibmResp)

    // Check if we have results
    if len(ibmResp.Results) == 0 {
        return nil, fmt.Errorf("no results returned from IBM model")
    }

    // Get predictions directly
    preds := ibmResp.Results[0].Predictions
    return i.convertToResult(preds.Toxic, preds.SevereToxic, preds.Obscene, preds.Threat, preds.Insult, preds.IdentityHate, content), nil
}

// convertToResult converts IBM scores to our ToxicityResult format
func (i *IBMAnalyzer) convertToResult(toxic, severeToxic, obscene, threat, insult, identityHate float64, content string) *ToxicityResult {
    
    result := &ToxicityResult{
        Score:       toxic * 100,
        IsFlagged:   toxic > 0.5 || threat > 0.5,
        ToxicWords:  []string{},
        Suggestions: []string{},
        Categories:  []ToxicityCategory{},
    }

    // Add all toxicity categories
    result.Categories = append(result.Categories,
        ToxicityCategory{
            Name:        "Toxicity",
            Score:       toxic * 100,
            Detected:    toxic > 0.5,
            Description: "General toxic content",
        },
        ToxicityCategory{
            Name:        "Severe Toxicity",
            Score:       severeToxic * 100,
            Detected:    severeToxic > 0.5,
            Description: "Extremely toxic content",
        },
        ToxicityCategory{
            Name:        "Obscene",
            Score:       obscene * 100,
            Detected:    obscene > 0.5,
            Description: "Obscene or vulgar language",
        },
        ToxicityCategory{
            Name:        "Threat",
            Score:       threat * 100,
            Detected:    threat > 0.5,
            Description: "Threatening content",
        },
        ToxicityCategory{
            Name:        "Insult",
            Score:       insult * 100,
            Detected:    insult > 0.5,
            Description: "Insulting language",
        },
        ToxicityCategory{
            Name:        "Identity Hate",
            Score:       identityHate * 100,
            Detected:    identityHate > 0.5,
            Description: "Attacks based on identity",
        },
    )

    // Determine overall severity based on highest score
    highestScore := toxic
    if threat > highestScore {
        highestScore = threat
    }
    if severeToxic > highestScore {
        highestScore = severeToxic
    }
    
    switch {
    case highestScore > 0.8:
        result.Severity = "critical"
    case highestScore > 0.6:
        result.Severity = "high"
    case highestScore > 0.4:
        result.Severity = "medium"
    case highestScore > 0.2:
        result.Severity = "low"
    default:
        result.Severity = "none"
    }

    // Set confidence based on model
    result.Confidence = 0.9

    // Determine sentiment based on scores
    if toxic < 0.2 && insult < 0.2 && threat < 0.2 {
        result.Sentiment = "positive"
    } else if toxic < 0.4 {
        result.Sentiment = "neutral"
    } else {
        result.Sentiment = "negative"
    }

    // Add suggestions if content is flagged
    if result.IsFlagged {
        if threat > 0.5 {
            result.Suggestions = append(result.Suggestions,
                "⚠️ Threatening content detected. This violates our safety policy.")
        }
        if identityHate > 0.5 {
            result.Suggestions = append(result.Suggestions,
                "⚠️ Hate speech based on identity is strictly prohibited.")
        }
        if insult > 0.5 {
            result.Suggestions = append(result.Suggestions,
                "⚠️ Insulting language detected. Please communicate respectfully.")
        }
        if obscene > 0.5 {
            result.Suggestions = append(result.Suggestions,
                "⚠️ Obscene language is not allowed.")
        }
        result.Suggestions = append(result.Suggestions,
            "This content has been flagged for review.")
    }

    return result
}

// CheckHealth checks if IBM service is available
func (i *IBMAnalyzer) CheckHealth() error {
    resp, err := i.client.Get("http://localhost:5000/model/predict")
    if err != nil {
        return fmt.Errorf("cannot connect to IBM model: %v", err)
    }
    defer resp.Body.Close()
    
    return nil
}

// Global IBM analyzer instance
var GlobalIBMAnalyzer = NewIBMAnalyzer()

// AnalyzeToxicityWithIBM is the main export function
func AnalyzeToxicityWithIBM(content string) (*ToxicityResult, error) {
    return GlobalIBMAnalyzer.Analyze(content)
}
