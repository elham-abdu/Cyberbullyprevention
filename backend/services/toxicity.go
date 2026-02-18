package services

import (
    "strings"
    "log"
)

func AnalyzeToxicity(content string) (float64, bool, error) {
    // Simple built-in toxicity detection
    toxicWords := []string{
        "bad", "hate", "stupid", "idiot", "dumb", "ugly", 
        "kill", "die", "hell", "damn", "shit", "fuck",
        "loser", "worthless", "trash", "garbage",
    }

    text := strings.ToLower(content)
    toxicCount := 0
    
    // Count toxic words
    for _, word := range toxicWords {
        if strings.Contains(text, word) {
            toxicCount++
            log.Printf("Found toxic word: %s", word)
        }
    }

    // Calculate score (0-100)
    var score float64
    if toxicCount > 0 {
        // Each toxic word adds 20% to the score, max 100%
        score = float64(toxicCount * 20)
        if score > 100 {
            score = 100
        }
    } else {
        // No toxic words found, give a low random-ish score
        score = 10
    }

    // Flag if score is 50 or above
    isFlagged := score >= 50

    log.Printf("Toxicity analysis - Score: %.2f%%, Flagged: %v, Toxic words found: %d", 
        score, isFlagged, toxicCount)

    return score, isFlagged, nil
}
