# CyberGuard API

A production-ready cyberbullying prevention API with ML-powered toxicity detection.

## 🚀 Overview

CyberGuard provides real-time content moderation through a REST API, using IBM's pre-trained BERT model to detect toxic language across 6 categories: toxicity, severe toxicity, obscenity, threats, insults, and identity hate.

## ✨ Features

- **JWT Authentication** - Secure user registration/login with role-based access
- **ML Toxicity Detection** - Real-time content analysis with 90%+ accuracy
- **Post Management** - Create, read, update, delete with automatic screening
- **Admin Dashboard** - Review and manage flagged content
- **Detailed Analysis** - Multi-category toxicity scores with severity levels

## 🛠 Tech Stack

**Backend:** Go 1.21+, GORM, PostgreSQL, JWT, Gorilla Mux  
**ML Engine:** IBM MAX Toxic Comment Classifier (BERT) in Docker  
**Frontend:** React 18, TypeScript, Vite, TailwindCSS, React Router

## 📦 Installation

### Prerequisites
- Go 1.21+, Docker, PostgreSQL 14+, Node.js 18+

### Quick Start

```bash
# Clone and setup
git clone https://github.com/yourusername/cyberguard-api.git
cd cyberguard-api

# Start ML model (keep this terminal open)
docker run -it -p 5000:5000 codait/max-toxic-comment-classifier

# Backend (new terminal)
cd backend
cp .env.example .env  # Add your DB credentials
go mod download
go run main.go

# Frontend (new terminal)
cd frontend
npm install
npm run dev
```

## 🔌 API Endpoints

### Public
- `POST /register` - Create account
- `POST /login` - Authenticate & get JWT

### Protected (JWT required)
- `GET /me` - Current user info
- `GET /me/posts` - User's posts
- `POST /me/posts/create` - Create post with toxicity analysis
- `PUT /me/posts/edit` - Update post
- `DELETE /me/posts/delete` - Remove post

### Admin (JWT + admin role)
- `GET /admin/flagged-posts` - View flagged content
- `POST /admin/posts/mark-safe` - Approve content
- `DELETE /admin/posts/delete-flagged` - Remove toxic content

## 🧠 ML Analysis Example

```json
Input: "I'm going to kill you"
Response: {
  "toxic": 0.94,
  "threat": 0.96,
  "insult": 0.10,
  "severe_toxic": 0.22
}
Result: FLAGGED (Critical severity)
```

## ⚙️ Configuration

### Backend (.env)
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=cyberguard
DB_PORT=5432
JWT_SECRET=your-secret-key
```

### Frontend (.env)
```env
VITE_API_URL=http://localhost:8080
```

## 📊 Database Schema

```sql
users: id, email, password_hash, role, timestamps
posts: id, user_id, content, toxicity_score, is_flagged, severity, sentiment, timestamps
```

## 🔒 Security

- bcrypt password hashing
- 24-hour JWT expiration
- Role-based access control
- CORS configured for frontend
- Parameterized queries via GORM

## 🐳 Docker Deployment

```bash
# Build and run all services
docker-compose up -d
```

## 📈 Performance

- Response time: 200-500ms per request
- Throughput: 50+ req/sec
- ML accuracy: ~90% on benchmarks

## 📝 License

MIT © Elham Abdu

---

**Made for a safer internet**
