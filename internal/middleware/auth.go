package middleware

import (
	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func TelegramAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		initData := c.GetHeader("tma")
		if initData == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
		telegramBotToken := "6400163949:AAGz5hrq3L_176NvCSeLM4tPrxQsCSJzdUg"

		err := initdata.Validate(initData, telegramBotToken, 0)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid initData"})
			return
		}

		value, err := initdata.Parse(initData)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized/Invalid initData ERRORRRRRR"})
			return
		}
		c.Set("tma", value.User)
		c.Next()
	}
}
