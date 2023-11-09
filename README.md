# Online Offline Indicator
1. docker-compose up -d
2. 4. create new db `test`
    1. create new table `pulse`
2. `go get -u github.com/gin-gonic/gin`
5. `go get github.com/go-sql-driver/mysql`
6. `go run main.go`
7. curl --location 'http://localhost:9000/heartbeats' \
--header 'Content-Type: application/json' \
--data '{"user_id":1}'
8. `curl --location 'http://localhost:9000/heartbeats/status/1'`
4. `curl -X GET http://localhost:9000/heartbeats/status_nonop\?user_ids\=1,2,3,4,5`
5. `curl -X GET http://localhost:9000/heartbeats/status\?user_ids\=1,2,3,4,5`
