// cmd/server/main.go
package main

import (
	"log"
	"math/rand"
	"monopoly/internal/api/handler"
	"monopoly/internal/manager"
	"monopoly/internal/user"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// 初始化依赖
	userManager := user.NewManager()
	gameManager := manager.NewGameManager()

	// 初始化处理器
	userHandler := handler.NewUserHandler(userManager)
	gameHandler := handler.NewGameHandler(gameManager, userManager)

	// 创建路由器
	r := mux.NewRouter()

	// API 路由
	apiRouter := r.PathPrefix("/api").Subrouter()

	// 用户相关路由
	apiRouter.HandleFunc("/users", userHandler.Create).Methods("POST")
	apiRouter.HandleFunc("/users/{userId}", userHandler.Get).Methods("GET")
	apiRouter.HandleFunc("/users/{userId}", userHandler.Update).Methods("PUT")
	apiRouter.HandleFunc("/users/{userId}", userHandler.Delete).Methods("DELETE")
	apiRouter.HandleFunc("/users/{userId}/coins", userHandler.AddCoins).Methods("POST")
	apiRouter.HandleFunc("/users/{userId}/coins/deduct", userHandler.DeductCoins).Methods("POST")
	apiRouter.HandleFunc("/users/{userId}/transactions", userHandler.GetUserTransactions).Methods("GET")
	apiRouter.HandleFunc("/users/{userId}/games", userHandler.GetUserGames).Methods("GET")
	apiRouter.HandleFunc("/users/{userId}/balance", userHandler.CheckUserBalance).Methods("GET")

	// 游戏相关路由
	apiRouter.HandleFunc("/games", gameHandler.Create).Methods("POST")
	apiRouter.HandleFunc("/games/{gameId}", gameHandler.Get).Methods("GET")
	apiRouter.HandleFunc("/games/{gameId}/join", gameHandler.Join).Methods("POST")
	apiRouter.HandleFunc("/games/{gameId}/start", gameHandler.StartGame).Methods("POST")
	apiRouter.HandleFunc("/games/{gameId}/roll", gameHandler.RollDice).Methods("POST")
	apiRouter.HandleFunc("/games/{gameId}/properties/{position}/buy", gameHandler.BuyProperty).Methods("POST")
	apiRouter.HandleFunc("/games/{gameId}/properties/{position}/upgrade", gameHandler.UpgradeProperty).Methods("POST")
	apiRouter.HandleFunc("/games/{gameId}/end-turn", gameHandler.EndTurn).Methods("POST")
	apiRouter.HandleFunc("/games/{gameId}/status", gameHandler.GetGameStatus).Methods("GET")
	apiRouter.HandleFunc("/games/{gameId}/players/{playerId}", gameHandler.GetPlayerStatus).Methods("GET")
	apiRouter.HandleFunc("/games/{gameId}/players/{playerId}/leave", gameHandler.LeaveGame).Methods("POST")

	// 中间件
	apiRouter.Use(loggingMiddleware)
	apiRouter.Use(recoveryMiddleware)

	// 服务器配置
	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	log.Printf("Server starting on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

// loggingMiddleware 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// recoveryMiddleware 恢复中间件
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
