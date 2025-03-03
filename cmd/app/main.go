package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/rs/cors"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/handlers"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/middleware"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/pb"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/repositories"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/services"
	gapi_services "github.com/sonnyvictok/miniapp_taptoearn/internal/services/gapi"

	worker "github.com/sonnyvictok/miniapp_taptoearn/internal/workers"
	"github.com/sonnyvictok/miniapp_taptoearn/pkg/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	r := gin.New()

	dbPostgres, err := database.NewPostgresDB()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{
	// 		"http://localhost:3000",
	// 		"https://4706-2405-4803-c86c-5730-75da-d2ac-ffc3-781d.ngrok-free.app",
	// 		"https://2451-2405-4803-c86c-5730-75da-d2ac-ffc3-781d.ngrok-free.app",
	// 		"*",
	// 	},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"tma", "Content-Type"},
	// 	AllowCredentials: true,
	// }))

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

	// go runGrpcServer(userRepo, &taskDistributor)
	go runGRPCServerGateway(userRepo, redisConnOpt, taskDistributor)

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

func runGRPCServerGateway(userRepo *repositories.UserRepository, redisConnOpt asynq.RedisClientOpt, taskDistributor worker.TaskDistributor) {
	serverGapi, err := gapi_services.NewServerGapi(userRepo, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc server")
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	c := cors.New(
		cors.Options{
			AllowedOrigins: []string{
				"http://localhost:3000",
				"https://4706-2405-4803-c86c-5730-75da-d2ac-ffc3-781d.ngrok-free.app",
				"https://2451-2405-4803-c86c-5730-75da-d2ac-ffc3-781d.ngrok-free.app",
				"*",
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"tma", "Content-Type"},
			AllowCredentials: true,
		},
	)
	headerMatcher := runtime.WithIncomingHeaderMatcher(customHeaderMatcher)
	grpcMux := runtime.NewServeMux(jsonOption, headerMatcher)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := pb.RegisterTapServiceHandlerServer(ctx, grpcMux, serverGapi); err != nil {
		log.Fatal().Err(err).Msg("failed to register handler server")
	}
	server := &http.Server{
		Addr:    ":9090",
		Handler: c.Handler(grpcMux),
	}

	// go runTaskProcessor(redisConnOpt, serverGapi.userService)
	// Start server
	log.Info().Msgf("Starting HTTP gateway server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}

func customHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "tma", "authorization":
		return key, true
	default:
		return "", false
	}
}

func runGrpcServer(userRepo *repositories.UserRepository, taskDistributor worker.TaskDistributor) {
	serverGapi, err := gapi_services.NewServerGapi(userRepo, taskDistributor)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc server")
	}

	fmt.Println(serverGapi)

	gprcLogger := grpc.UnaryInterceptor(gapi_services.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	pb.RegisterTapServiceServer(grpcServer, serverGapi)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
