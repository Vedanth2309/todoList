# Personal Tracker
React + Vite + Tailwind + Recharts frontend, Go + Fiber API, MongoDB Atlas (source of truth), Elasticsearch (search).

## Setup
Prereqs: Node 18+, Go 1.22+, a MongoDB Atlas cluster, Elasticsearch 8 (e.g. `docker run -p 9200:9200 -e discovery.type=single-node -e xpack.security.enabled=false elasticsearch:8.14.0`).

    cp .env.example .env      # fill MONGODB_URI, JWT_SECRET, ELASTICSEARCH_URL
    cd backend && go mod tidy && go run .      # :8080
    cd frontend && npm install && npm run dev  # :5173 (proxies /api)

## API
All routes except register/login need `Authorization: Bearer <jwt>`; user ID comes only from the JWT.
`POST /api/auth/register|login`, `GET /api/auth/me`, CRUD on `/api/{tasks,habits,events,reminders,diary,notes}`,
`POST /api/habits/:id/checkin`, `GET /api/habits/:id/checkins`, `GET /api/search?q=&type=&tag=`,
`GET /api/analytics/dashboard`, `GET|PUT /api/settings`.
Responses: `{success, data}` or `{success:false, message}`.
