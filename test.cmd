#docker-compose down
#docker-compose up -d
#docker-compose ps

go run cmd/api/main.go

curl http://localhost:8080/api/health

curl -X POST http://localhost:8080/api/register -H "Content-Type: application/json" -d "{\"username\": \"tester_001\", \"password\": \"tester_pass123\", \"email\": \"tester_001@example.com\"}"

curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d "{\"email\": \"tester_001@example.com\", \"password\": \"tester_pass123\"}"

curl -X POST http://localhost:8080/api/posts -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3Rlcl8wMDFAZXhhbXBsZS5jb20iLCJ1c2VybmFtZSI6InRlc3Rlcl8wMDEiLCJleHAiOjE3ODMwMDM0ODMsImlhdCI6MTc4MjkxNzA4M30.8dhhroFGOzTsOMOom0yq27qh7yLU5a9sl7B-olxcrH4"  -d "{\"title\": \"Test post 1A\", \"content\": \"test  on post 1.\"}"
