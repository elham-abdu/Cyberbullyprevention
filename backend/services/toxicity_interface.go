package services

import (
    "strings"
)

// ToxicityResult contains detailed analysis results
type ToxicityResult struct {
    Score       float64           `json:"score"`
    IsFlagged   bool              `json:"is_flagged"`
    Severity    string            `json:"severity"`
    Categories  []ToxicityCategory `json:"categories"`
    ToxicWords  []string          `json:"toxic_words"`
    Sentiment   string            `json:"sentiment"`
    Confidence  float64           `json:"confidence"`
    Suggestions []string          `json:"suggestions"`
}

// ToxicityCategory represents different types of toxic content
type ToxicityCategory struct {
    Name        string  `json:"name"`
    Score       float64 `json:"score"`
    Detected    bool    `json:"detected"`
    Description string  `json:"description"`
}

// ToxicityService defines the interface for toxicity analysis
type ToxicityService interface {
    Analyze(content string) (*ToxicityResult, error)
}

// GetToxicityService returns the appropriate toxicity service
func GetToxicityService() ToxicityService {
    // Create a new IBM analyzer to check health
    ibmAnalyzer := NewIBMAnalyzer()
    
    // Try IBM first
    if err := ibmAnalyzer.CheckHealth(); err == nil {
        return GlobalIBMAnalyzer
    }
    
    // Fallback to rule-based if IBM not available
    return NewRuleBasedAnalyzer()
}

// RuleBasedAnalyzer as fallback
type RuleBasedAnalyzer struct {}

func NewRuleBasedAnalyzer() *RuleBasedAnalyzer {
    return &RuleBasedAnalyzer{}
}

func (r *RuleBasedAnalyzer) Analyze(content string) (*ToxicityResult, error) {
    // Use rule-based analysis
    return analyzeWithRules(content), nil
}

// analyzeWithRules provides rule-based fallback
func analyzeWithRules(content string) *ToxicityResult {
    result := &ToxicityResult{
        Categories:  []ToxicityCategory{},
        ToxicWords:  []string{},
        Suggestions: []string{},
    }

    text := strings.ToLower(content)
    words := strings.Fields(text)

    // Define word lists
    profanityList := []string{"fuck", "shit", "damn", "hell", "ass", "bitch"}
    insultList := []string{"stupid", "idiot", "dumb", "moron", "loser", "dummy"}
    threatList := []string{"kill", "die", "hurt", "attack", "destroy", "beat"}
    hateSpeechList := []string{"hate", "racist", "sexist", "nazi", "discriminate"}

    // Count occurrences in words
    profanityCount := 0
    insultCount := 0
    threatCount := 0
    hateCount := 0

    // Check each word individually
    for _, word := range words {
        wordLower := strings.ToLower(word)
        
        // Check profanity
        for _, badWord := range profanityList {
            if strings.Contains(wordLower, badWord) {
                profanityCount++
                result.ToxicWords = append(result.ToxicWords, word)
                break
            }
        }
        
        // Check insults
        for _, badWord := range insultList {
            if strings.Contains(wordLower, badWord) {
                insultCount++
                result.ToxicWords = append(result.ToxicWords, word)
                break
            }
        }
        
        // Check threats
        for _, badWord := range threatList {
            if strings.Contains(wordLower, badWord) {
                threatCount++
                result.ToxicWords = append(result.ToxicWords, word)
                break
            }
        }
        
        // Check hate speech
        for _, badWord := range hateSpeechList {
            if strings.Contains(wordLower, badWord) {
                hateCount++
                result.ToxicWords = append(result.ToxicWords, word)
                break
            }
        }
    }

    // Calculate scores (0-100)
    profanityScore := float64(profanityCount * 20)
    insultScore := float64(insultCount * 20)
    threatScore := float64(threatCount * 25)
    hateScore := float64(hateCount * 30)

    // Cap at 100
    if profanityScore > 100 { profanityScore = 100 }
    if insultScore > 100 { insultScore = 100 }
    if threatScore > 100 { threatScore = 100 }
    if hateScore > 100 { hateScore = 100 }

    // Add categories
    result.Categories = append(result.Categories,
        ToxicityCategory{Name: "Profanity", Score: profanityScore, 
            Detected: profanityCount > 0, Description: "Profane language"},
        ToxicityCategory{Name: "Insults", Score: insultScore, 
            Detected: insultCount > 0, Description: "Insulting language"},
        ToxicityCategory{Name: "Threats", Score: threatScore, 
            Detected: threatCount > 0, Description: "Threatening language"},
        ToxicityCategory{Name: "Hate Speech", Score: hateScore, 
            Detected: hateCount > 0, Description: "Hate speech"},
    )

    // Calculate overall score
    totalScore := (profanityScore*0.2 + insultScore*0.25 + threatScore*0.3 + hateScore*0.25)

    result.Score = totalScore
    result.IsFlagged = totalScore >= 30
    result.Severity = getSeverity(totalScore)
    result.Sentiment = getSentiment(text)
    result.Confidence = 0.6 + (float64(len(result.ToxicWords)) * 0.1)
    if result.Confidence > 0.95 {
        result.Confidence = 0.95
    }

    return result
}

// Helper functions
func getSeverity(score float64) string {
    switch {
    case score >= 70:
        return "high"
    case score >= 40:
        return "medium"
    case score >= 20:
        return "low"
    default:
        return "none"
    }
}

func getSentiment(text string) string {
    positiveWords := []string{"good", "great", "awesome", "excellent", "love", "thanks", "perfect"}
    negativeWords := []string{"bad", "hate", "awful", "terrible", "worst", "horrible"}

    posCount := 0
    negCount := 0

    for _, word := range positiveWords {
        if strings.Contains(text, word) {
            posCount++
        }
    }

    for _, word := range negativeWords {
        if strings.Contains(text, word) {
            negCount++
        }
    }

    if posCount > negCount*2 {
        return "positive"
    } else if negCount > posCount*2 {
        return "negative"
    } else if posCount == 0 && negCount == 0 {
        return "neutral"
    } else {
        return "mixed"
    }
}
