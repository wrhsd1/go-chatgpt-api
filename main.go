package main

import (
	"log"
	"os"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"

	"github.com/maxduke/go-chatgpt-api/api"
	"github.com/maxduke/go-chatgpt-api/api/chatgpt"
	"github.com/maxduke/go-chatgpt-api/api/imitate"
	"github.com/maxduke/go-chatgpt-api/api/platform"
	_ "github.com/maxduke/go-chatgpt-api/env"
	"github.com/maxduke/go-chatgpt-api/middleware"
)

func init() {
	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)
}

func main() {
    router := gin.Default()

    router.Use(middleware.CORS())
    router.Use(middleware.Authorization())

    path_key := os.Getenv("PATH_KEY")
    passwordGroup := router.Group("/" + path_key)
    {
        setupChatGPTAPIs(passwordGroup)
        setupPlatformAPIs(passwordGroup)
        setupImitateAPIs(passwordGroup)
    }
    setupPandoraAPIs(router, path_key)
    router.NoRoute(api.Proxy)

    router.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, api.ReadyHint)
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    err := router.Run(":" + port)
    if err != nil {
        log.Fatal("failed to start server: " + err.Error())
    }
}

func setupChatGPTAPIs(router *gin.RouterGroup) {
    chatgptGroup := router.Group("/chatgpt")
	{
		chatgptGroup.POST("/login", chatgpt.Login)
		chatgptGroup.POST("/backend-api/login", chatgpt.Login) // add support for other projects

		conversationGroup := chatgptGroup.Group("/backend-api/conversation")
		{
			conversationGroup.POST("", chatgpt.CreateConversation)
		}
	}
}

func setupPlatformAPIs(router *gin.RouterGroup) {
    platformGroup := router.Group("/platform")
	{
		platformGroup.POST("/login", platform.Login)
		platformGroup.POST("/v1/login", platform.Login)

		apiGroup := platformGroup.Group("/v1")
		{
			apiGroup.POST("/chat/completions", platform.CreateChatCompletions)
			apiGroup.POST("/completions", platform.CreateCompletions)
		}
	}
}

func setupPandoraAPIs(router *gin.Engine, path_key string) {
    router.Any("/api/*path", func(c *gin.Context) {
        c.Request.URL.Path = strings.ReplaceAll(c.Request.URL.Path, "/api", "/" + path_key + "/chatgpt/backend-api")
        router.HandleContext(c)
    })
}

func setupImitateAPIs(router *gin.RouterGroup) {
    imitateGroup := router.Group("/imitate")
	{
		imitateGroup.POST("/login", chatgpt.Login)

		apiGroup := imitateGroup.Group("/v1")
		{
			apiGroup.POST("/chat/completions", imitate.CreateChatCompletions)
		}
	}
}
