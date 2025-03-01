package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/handlers"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/middleware"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/repositories"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/services"
	worker "github.com/sonnyvictok/miniapp_taptoearn/internal/workers"
	"github.com/sonnyvictok/miniapp_taptoearn/pkg/database"
)

func main() {
	r := gin.New()

	dbPostgres, err := database.NewPostgresDB()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://c02b-2405-4803-c844-6a0-a703-9d2b-c2af-b8c4.ngrok-free.app"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"tma", "Content-Type"},
		AllowCredentials: true,
	}))

	// defer dbPostgres.Close()
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0,
	// })
	// ctx := context.Background()
	// if err := redisClient.Ping(ctx).Err(); err != nil {
	// 	log.Fatalf("Failed to connect to Redis: %v", err)
	// }
	// defer redisClient.Close()
	// create redis connection
	redisConnOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
		DB:   0,
	}

	userRepo := repositories.NewUserRepository(dbPostgres)
	userServices := services.NewUserService(userRepo)
	taskDistributor := worker.NewRedisTaskDistributor(redisConnOpt)

	userHandler := handlers.NewUserHandler(userServices, taskDistributor)
	go runTaskProcessor(redisConnOpt, userServices)

	r.Use(middleware.TelegramAuthMiddleware())
	r.POST("/clicktoearn", userHandler.ClickToEarn)
	r.GET("/user", userHandler.GetUser)
	r.POST("/createuser", userHandler.CreateUser)

	r.Run(":8080")

	// indata := "user=%7B%22id%22%3A5204989815%2C%22first_name%22%3A%22Sonny%22%2C%22last_name%22%3A%22Victok%22%2C%22username%22%3A%22SonnyVitok%22%2C%22language_code%22%3A%22en%22%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2FrUi7bmeMCZfIcid5iCSQZWMRUTPriPLJVNmdvts9Kc1umIhBNJ2c4_lJlVCeSlX0.svg%22%7D&chat_instance=-4534184936334365392&chat_type=private&auth_date=1740491181&signature=MwgViiMS0tIpe5P1U1BZRRxqnkZ9VLeKU_5Nbod5ItuCtzAHnr-kBK7UPaHj6iJSyLF_YIN7XAVJ22rMC-1OAg&hash=2c5b7e01eac3cae892420f25010702e76681d2dc848569278606759a388c5c26"

	// // Your secret bot token.
	// token := "6400163949:AAGz5hrq3L_176NvCSeLM4tPrxQsCSJzdUg"
}

func runTaskProcessor(redisConnOpt asynq.RedisClientOpt, userServices *services.UserService) {
	taskProcessor := worker.NewRedisTaskProcessor(redisConnOpt, userServices)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")

	}
}

// user=%7B%22id%22%3A5204989815%2C%22first_name%22%3A%22Sonny%22%2C%22last_name%22%3A%22Victok%22%2C%22username%22%3A%22SonnyVitok%22%2C%22language_code%22%3A%22en%22%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2FrUi7bmeMCZfIcid5iCSQZWMRUTPriPLJVNmdvts9Kc1umIhBNJ2c4_lJlVCeSlX0.svg%22%7D&chat_instance=-4534184936334365392&chat_type=private&auth_date=1740494970&signature=jszRckMd4AZn0o2OeYwzRuHfiKuAnUzwKUVtU_SOi-XhN3-Ml3fRu3TvvMmhx_PHqUmkfKktQcUmtyMH0u_tBQ&hash=17b32d1151055f447f3b75484bb48f753570ccd3bcd8f16fe46f20544a29483e

// user=%7B%22id%22%3A6183423690%2C%22first_name%22%3A%22Corle%F0%9F%A6%B4%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22Trysuccesss%22%2C%22language_code%22%3A%22en%22%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2FH-scBa0O4wzmjoeN7-jKuq_L2fNm8EuvwsCKwfKgMy-aqdE1P99Fx3ay6tmbl439.svg%22%7D&chat_instance=3794538746268988861&chat_type=private&auth_date=1740771160&signature=d56dXT9IYgubyvPQGqMOMPr48d6_VwiTExnwSjoqTtj8LH-7JUx55bhobwnpFWleJcH1BU5K2v1ZAt3GGSwtCw&hash=dd1181679b66aa067cae20127c423bedb58725559297e2130b11a53b3aee5dde
