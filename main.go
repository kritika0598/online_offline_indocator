package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "secret"
	hostname = "127.0.0.1:3306"
	dbname   = "test"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

var DB *sql.DB

func init() {
	_db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err)
	}
	DB = _db
}

func main() {
	ge := gin.Default()

	ge.POST("/heartbeats", func(ctx *gin.Context) {
		data := map[string]interface{}{}
		ctx.Bind(&data)
		_, err := DB.Exec("REPLACE INTO pulse (user_id, last_hb) VALUES (?, ?);", data["user_id"], time.Now().Unix())
		if err != nil {
			panic(err)
		}
		ctx.JSON(200, map[string]interface{}{"message": "ok"})
	})

	ge.GET("/heartbeats/status/:user_id", func(ctx *gin.Context) {
		var lastHb int
		row := DB.QueryRow("SELECT last_hb FROM pulse WHERE user_id = ?;", ctx.Param("user_id"))
		row.Scan(&lastHb)
		ctx.JSON(200, map[string]interface{}{"is_online": lastHb > int(time.Now().Unix())-30})
	})

	// BATCH
	ge.GET("/heartbeats/status", func(ctx *gin.Context) {
		rows, err := DB.Query("SELECT user_id, last_hb FROM pulse WHERE user_id IN (?);", ctx.Query("user_ids"))
		if err != nil {
			panic(err)
		}
		var statusMap = make(map[string]bool)
		var userId, lastHb int
		for rows.Next() {
			if err := rows.Scan(&userId, &lastHb); err != nil {
				panic(err)
			}
			statusMap[fmt.Sprintf("%d", userId)] = lastHb > int(time.Now().Unix())-30
		}
		rows.Close()
		ctx.JSON(200, statusMap)
	})

	// NON Optimised
	ge.GET("/heartbeats/status_nonop", func(ctx *gin.Context) {
		var statusMap = make(map[string]bool)
		for _, userID := range strings.Split(ctx.Query("user_ids"), ",") {
			var lastHb int
			row := DB.QueryRow("SELECT last_hb FROM pulse WHERE user_id = ?;", userID)
			row.Scan(&lastHb)
			statusMap[userID] = lastHb > int(time.Now().Unix())-30
		}
		ctx.JSON(200, statusMap)
	})

	ge.Run(":9000")
}
