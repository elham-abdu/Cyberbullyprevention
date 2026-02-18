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

// IBMResponse matches the IBM MAX model response
type IBMResponse struct {
    Status  string        `json:"status"`
    Results []IBMResult   `json:"results"`
}

// IBMResult contains toxicity scores for one text
type IBMResult struct {
    ToxicityScore     float64 `json:"toxicity"`
    SevereToxicity    float64 `json:"severe_toxicity"`
    Obscene           float64 `json:"obscene"`
    Threat            float64 `json:"threat"`
    Insult            float64 `json:"insult"`
    IdentityAttack    float64 `json:"identity_attack"`
    SexualExplicit    float64 `json:"sexual_explicit"`
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

    // Log request (for debugging)
    log.Printf("Sending request to IBM model: %s", string(jsonData))

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

    // Log response (for debugging)
    log.Printf("Received response from IBM model: %+v", ibmResp)

    // Check if we have results
    if len(ibmResp.Results) == 0 {
        return nil, fmt.Errorf("no results returned from IBM model")
    }

    // Convert IBM result to our ToxicityResult format
    return i.convertToResult(ibmResp.Results[0], content), nil
}

// convertToResult converts IBM result to our ToxicityResult format
func (i *IBMAnalyzer) convertToResult(ibmResult IBMResult, content string) *ToxicityResult {
    result := &ToxicityResult{
        Score:       ibmResult.ToxicityScore * 100,
        IsFlagged:   ibmResult.ToxicityScore > 0.5,
        ToxicWords:  []string{},
        Suggestions: []string{},
        Categories:  []ToxicityCategory{},
    }

    // Add all toxicity categories
    result.Categories = append(result.Categories,
        ToxicityCategory{
            Name:        "Toxicity",
            Score:       ibmResult.ToxicityScore * 100,
            Detected:    ibmResult.ToxicityScore > 0.5,
            Description: "General toxic content",
        },
        ToxicityCategory{
            Name:        "Severe Toxicity",
            Score:       ibmResult.SevereToxicity * 100,
            Detected:    ibmResult.SevereToxicity > 0.5,
            Description: "Extremely toxic content",
        },
        ToxicityCategory{
            Name:        "Obscene",
            Score:       ibmResult.Obscene * 100,
            Detected:    ibmResult.Obscene > 0.5,
            Description: "Obscene or vulgar language",
        },
        ToxicityCategory{
            Name:        "Threat",
            Score:       ibmResult.Threat * 100,
            Detected:    ibmResult.Threat > 0.5,
            Description: "Threatening content",
        },
        ToxicityCategory{
            Name:        "Insult",
            Score:       ibmResult.Insult * 100,
            Detected:    ibmResult.Insult > 0.5,
            Description: "Insulting language",
        },
        ToxicityCategory{
            Name:        "Identity Attack",
            Score:       ibmResult.IdentityAttack * 100,
            Detected:    ibmResult.IdentityAttack > 0.5,
            Description: "Attacks based on identity",
        },
        ToxicityCategory{
            Name:        "Sexual Explicit",
            Score:       ibmResult.SexualExplicit * 100,
            Detected:    ibmResult.SexualExplicit > 0.5,
            Description: "Sexually explicit content",
        },
    )

    // Determine overall severity based on highest score
    highestScore := ibmResult.ToxicityScore
    if ibmResult.SevereToxicity > highestScore {
        highestScore = ibmResult.SevereToxicity
    }
    if ibmResult.Threat > highestScore {
        highestScore = ibmResult.Threat
    }
    if ibmResult.IdentityAttack > highestScore {
        highestScore = ibmResult.IdentityAttack
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

    // Set confidence based on model (IBM model is quite accurate)
    result.Confidence = 0.9

    // Determine sentiment based on scores
    if ibmResult.ToxicityScore < 0.2 && ibmResult.Insult < 0.2 && ibmResult.Threat < 0.2 {
        result.Sentiment = "positive"
    } else if ibmResult.ToxicityScore < 0.4 {
        result.Sentiment = "neutral"
    } else {
        result.Sentiment = "negative"
    }

    // Add suggestions if content is flagged
    if result.IsFlagged {
        if ibmResult.Threat > 0.5 {
            result.Suggestions = append(result.Suggestions, 
                "⚠️ Threatening content detected. This violates our safety policy and may be reported to authorities.")
        }
        if ibmResult.IdentityAttack > 0.5 {
            result.Suggestions = append(result.Suggestions, 
                "⚠️ Hate speech based on identity is strictly prohibited.")
        }
        if ibmResult.Insult > 0.5 {
            result.Suggestions = append(result.Suggestions, 
                "⚠️ Insulting language detected. Please communicate respectfully.")
        }
        if ibmResult.SexualExplicit > 0.5 {
            result.Suggestions = append(result.Suggestions, 
                "⚠️ Sexually explicit content is not allowed.")
        }
        result.Suggestions = append(result.Suggestions, 
            "This content has been flagged for review by our moderation team.")
    }

    return result
}

// CheckHealth checks if IBM service is available
func (i *IBMAnalyzer) CheckHealth() error {
    // Use the client from the IBMAnalyzer struct
    resp, err := i.client.Get("http://localhost:5000/model/health")
    if err != nil {
        return fmt.Errorf("cannot connect to IBM model: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
    }
    
    return nil
}

// Global IBM analyzer instance
var GlobalIBMAnalyzer = NewIBMAnalyzer()

// AnalyzeToxicityWithIBM is the main export function
func AnalyzeToxicityWithIBM(content string) (*ToxicityResult, error) {
    return GlobalIBMAnalyzer.Analyze(content)
}
