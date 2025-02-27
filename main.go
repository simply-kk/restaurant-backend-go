package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "golang-restaurant-management/database"
    "golang-restaurant-management/middleware"
    "golang-restaurant-management/routes"

    "github.com/gin-gonic/gin"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    client, err := database.DbInstance()
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    if client == nil {
        log.Fatalf("MongoDB client is nil despite no error from DbInstance")
    }

    database.InitCollections(client)

    router := gin.Default()
    routes.UserRoutes(router)
    router.Use(middleware.Authentication())
    routes.FoodRoutes(router)
    routes.MenuRoutes(router)
    routes.TableRoutes(router)
    routes.OrderRoutes(router)
    routes.OrderItemRoutes(router)
    routes.InvoiceRoutes(router)

    go func() {
        fmt.Println("Server running on port:", port)
        if err := router.Run(":" + port); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit
    fmt.Println("\nShutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if client != nil {
        if err := client.Disconnect(ctx); err != nil {
            log.Println("Error disconnecting MongoDB:", err)
        } else {
            fmt.Println("MongoDB connection closed successfully.")
        }
    }

    fmt.Println("Server shutdown completed.")
}