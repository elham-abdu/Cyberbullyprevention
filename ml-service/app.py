from flask import Flask, request, jsonify
from flask_cors import CORS
import torch
from transformers import AutoTokenizer, AutoModelForSequenceClassification
import nltk
from nltk.sentiment import SentimentIntensityAnalyzer
import re

app = Flask(__name__)
CORS(app)

# Download NLTK data
nltk.download('vader_lexicon', quiet=True)
sia = SentimentIntensityAnalyzer()

# Load pre-trained toxicity model (using smaller model for demo)
try:
    tokenizer = AutoTokenizer.from_pretrained("unitary/toxic-bert")
    model = AutoModelForSequenceClassification.from_pretrained("unitary/toxic-bert")
    print("✅ Loaded toxicity model")
except:
    print("⚠️ Could not load toxicity model, using fallback")
    model = None
    tokenizer = None

# Toxicity categories
categories = [
    'toxicity', 'severe_toxicity', 'obscene', 'threat',
    'insult', 'identity_attack', 'sexual_explicit'
]

@app.route('/predict', methods=['POST', 'OPTIONS'])
def predict():
    if request.method == 'OPTIONS':
        return '', 200
    
    data = request.json
    text = data.get('text', '')
    
    # Use transformer model if available
    if model and tokenizer:
        inputs = tokenizer(text, return_tensors="pt", truncation=True, max_length=512)
        outputs = model(**inputs)
        predictions = torch.sigmoid(outputs.logits).detach().numpy()[0]
        
        result = {}
        for i, category in enumerate(categories):
            result[category] = float(predictions[i])
    else:
        # Fallback to rule-based with NLTK sentiment
        result = fallback_analysis(text)
    
    # Add sentiment analysis
    sentiment = sia.polarity_scores(text)
    result['sentiment_score'] = sentiment['compound']
    result['confidence'] = 0.85  # You could calculate this properly
    
    return jsonify(result)

def fallback_analysis(text):
    """Rule-based fallback when ML model isn't available"""
    text_lower = text.lower()
    
    # Toxic word lists
    toxic_words = {
        'toxicity': ['hate', 'stupid', 'idiot', 'dumb', 'ugly'],
        'severe_toxicity': ['kill', 'die', 'death', 'murder'],
        'insult': ['idiot', 'moron', 'retard', 'loser'],
        'threat': ['kill you', 'hurt you', 'destroy you'],
    }
    
    result = {}
    for category, words in toxic_words.items():
        score = sum(2 for word in words if word in text_lower)
        result[category] = min(score / 10, 1.0)
    
    # Add default values for missing categories
    result['obscene'] = result.get('toxicity', 0) * 0.7
    result['identity_attack'] = result.get('insult', 0) * 0.8
    result['sexual_explicit'] = 0.0
    
    return result

@app.route('/health', methods=['GET'])
def health():
    return jsonify({'status': 'healthy', 'model_loaded': model is not None})

if __name__ == '__main__':
    print("🚀 ML Service running on http://localhost:8000")
    app.run(host='0.0.0.0', port=8000, debug=True)
