# ML Service Setup Instructions

1. Install Python (3.8+)
2. Install dependencies:
   pip install -r requirements.txt

3. Run the service:
   python app.py

4. The service will be available at http://localhost:8000
   - POST /predict for toxicity analysis
   - GET /health for health check

Note: First run will download the ML model (~500MB)
