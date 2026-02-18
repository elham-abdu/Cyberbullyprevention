package services

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "regexp"
)

// MLToxicityAnalyzer uses ML models for toxicity detection
type MLToxicityAnalyzer struct {
    modelURL string // URL to your ML model service
    useLocal bool   // Whether to use local model or API
}

// MLPrediction represents the ML model output
type MLPrediction struct {
    ToxicityScore     float64            `json:"toxicity_score"`
    SevereToxicity    float64            `json:"severe_toxicity"`
    Insult            float64            `json:"insult"`
    Threat            float64            `json:"threat"`
    Profanity         float64            `json:"profanity"`
    IdentityAttack    float64            `json:"identity_attack"`
    SentimentScore    float64            `json:"sentiment_score"`
    Confidence        float64            `json:"confidence"`
    Categories        map[string]float64 `json:"categories"`
}

// REMOVED: ToxicityResult and ToxicityCategory type definitions
// They are now imported from toxicity_interface.go

// NewMLToxicityAnalyzer creates a new ML-based analyzer
func NewMLToxicityAnalyzer() *MLToxicityAnalyzer {
    return &MLToxicityAnalyzer{
        modelURL: "http://localhost:8000/predict",
        useLocal: true,
    }
}

// AnalyzeWithML performs ML-based toxicity analysis
func (m *MLToxicityAnalyzer) AnalyzeWithML(content string) (*ToxicityResult, error) {
    // Try ML model first
    mlResult, err := m.callMLModel(content)
    if err == nil && mlResult != nil {
        return m.convertMLResult(mlResult, content), nil
    }

    // Fallback to rule-based if ML fails
    log.Printf("ML model unavailable, falling back to rule-based: %v", err)
    return m.analyzeWithRules(content), nil
}

// callMLModel calls an external ML model service
func (m *MLToxicityAnalyzer) callMLModel(content string) (*MLPrediction, error) {
    requestBody := map[string]string{
        "text": content,
    }
    
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(m.modelURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var prediction MLPrediction
    err = json.NewDecoder(resp.Body).Decode(&prediction)
    if err != nil {
        return nil, err
    }

    return &prediction, nil
}

// analyzeWithRules provides rule-based fallback
func (m *MLToxicityAnalyzer) analyzeWithRules(content string) *ToxicityResult {
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

    // Also check for multi-word phrases
    fullTextLower := strings.ToLower(content)
    
    // Check for phrases in threats
    threatPhrases := []string{"kill you", "hurt you", "going to get you", "beat you up"}
    for _, phrase := range threatPhrases {
        if strings.Contains(fullTextLower, phrase) {
            threatCount += 2
            result.ToxicWords = append(result.ToxicWords, phrase)
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
    
    // Check for patterns (ALL CAPS, repetition)
    patternScore := m.detectPatterns(content)
    totalScore += patternScore

    // Word count analysis
    wordCount := len(words)
    if wordCount < 3 && totalScore > 50 {
        totalScore *= 0.8
    }

    if totalScore > 100 {
        totalScore = 100
    }

    result.Score = totalScore
    result.IsFlagged = totalScore >= 30
    result.Severity = m.getSeverity(totalScore)
    result.Sentiment = m.getSentimentFromWords(words, fullTextLower)
    result.Confidence = 0.6 + (float64(len(result.ToxicWords)) * 0.1)
    if result.Confidence > 0.95 {
        result.Confidence = 0.95
    }

    // Add suggestions based on what was found
    if result.IsFlagged {
        if threatCount > 0 {
            result.Suggestions = append(result.Suggestions, 
                "Threatening content detected. This violates our safety policy.")
        }
        if hateCount > 0 {
            result.Suggestions = append(result.Suggestions, 
                "Hate speech is not tolerated on this platform.")
        }
        if insultCount > 2 {
            result.Suggestions = append(result.Suggestions, 
                "Repeated insults detected. Please be respectful.")
        }
        result.Suggestions = append(result.Suggestions, 
            "Content has been flagged for review.")
    }

    return result
}

// getSentimentFromWords - more accurate sentiment analysis
func (m *MLToxicityAnalyzer) getSentimentFromWords(words []string, fullText string) string {
    positiveWords := []string{"good", "great", "awesome", "excellent", "love", "thanks", "perfect", "beautiful"}
    negativeWords := []string{"bad", "hate", "awful", "terrible", "worst", "horrible", "dislike"}

    posCount := 0
    negCount := 0

    // Check each word
    for _, word := range words {
        wordLower := strings.ToLower(word)
        for _, pos := range positiveWords {
            if wordLower == pos {
                posCount++
                break
            }
        }
        for _, neg := range negativeWords {
            if wordLower == neg {
                negCount++
                break
            }
        }
    }

    // Check for negations
    if strings.Contains(fullText, "not good") || strings.Contains(fullText, "not great") {
        negCount += 2
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

// detectPatterns detects suspicious patterns
func (m *MLToxicityAnalyzer) detectPatterns(content string) float64 {
    score := 0.0

    // Check for ALL CAPS
    capsCount := 0
    letterCount := 0
    for _, r := range content {
        if (r >= 'A' && r <= 'Z') {
            capsCount++
            letterCount++
        } else if (r >= 'a' && r <= 'z') {
            letterCount++
        }
    }
    if letterCount > 0 && float64(capsCount)/float64(letterCount) > 0.5 {
        score += 10
    }

    // Check for repeated characters
    repeatPattern := regexp.MustCompile(`(.)\1{3,}`)
    if repeatPattern.MatchString(content) {
        score += 10
    }

    // Check for excessive punctuation
    punctCount := strings.Count(content, "!") + strings.Count(content, "?")
    if punctCount > 3 {
        score += 5
    }

    return score
}

// getSeverity determines severity level
func (m *MLToxicityAnalyzer) getSeverity(score float64) string {
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

// convertMLResult converts ML model output to our result format
func (m *MLToxicityAnalyzer) convertMLResult(ml *MLPrediction, content string) *ToxicityResult {
    result := &ToxicityResult{
        Score:       ml.ToxicityScore * 100,
        IsFlagged:   ml.ToxicityScore > 0.4,
        Sentiment:   m.sentimentFromScore(ml.SentimentScore),
        Confidence:  ml.Confidence,
        ToxicWords:  []string{},
        Suggestions: []string{},
        Categories: []ToxicityCategory{
            {Name: "Severe Toxicity", Score: ml.SevereToxicity * 100, 
             Detected: ml.SevereToxicity > 0.3, Description: "Extremely toxic content"},
            {Name: "Insult", Score: ml.Insult * 100, 
             Detected: ml.Insult > 0.3, Description: "Insulting language"},
            {Name: "Threat", Score: ml.Threat * 100, 
             Detected: ml.Threat > 0.3, Description: "Threatening content"},
            {Name: "Profanity", Score: ml.Profanity * 100, 
             Detected: ml.Profanity > 0.3, Description: "Profane language"},
            {Name: "Identity Attack", Score: ml.IdentityAttack * 100, 
             Detected: ml.IdentityAttack > 0.3, Description: "Attacks on identity"},
        },
    }
    
    // Add custom categories from ML model
    for category, score := range ml.Categories {
        result.Categories = append(result.Categories, ToxicityCategory{
            Name: category, Score: score * 100, Detected: score > 0.3,
        })
    }
    
    result.Severity = m.determineSeverityML(result.Score, ml.Threat*100)
    
    return result
}

// determineSeverityML determines severity from ML results
func (m *MLToxicityAnalyzer) determineSeverityML(score, threatScore float64) string {
    if threatScore > 50 {
        return "critical"
    }
    switch {
    case score >= 70:
        return "high"
    case score >= 50:
        return "medium"
    case score >= 30:
        return "low"
    default:
        return "none"
    }
}

// sentimentFromScore converts ML sentiment score to string
func (m *MLToxicityAnalyzer) sentimentFromScore(score float64) string {
    switch {
    case score > 0.6:
        return "very_positive"
    case score > 0.2:
        return "positive"
    case score < -0.6:
        return "very_negative"
    case score < -0.2:
        return "negative"
    default:
        return "neutral"
    }
}

// Global ML analyzer instance
var GlobalMLAnalyzer = NewMLToxicityAnalyzer()

// AnalyzeToxicityML is the main export function
func AnalyzeToxicityML(content string) (*ToxicityResult, error) {
    return GlobalMLAnalyzer.AnalyzeWithML(content)
}
