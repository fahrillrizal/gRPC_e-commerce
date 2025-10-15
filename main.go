package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/fahrillrizal/ecommerce-grpc/internal/handler"
	"github.com/fahrillrizal/ecommerce-grpc/internal/repositories"
	"github.com/fahrillrizal/ecommerce-grpc/internal/services"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/pb/auth"
	"github.com/fahrillrizal/ecommerce-grpc/pb/product"
	"github.com/fahrillrizal/ecommerce-grpc/pkg/database"
	"github.com/fahrillrizal/ecommerce-grpc/pkg/middleware"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	godotenv.Load()

	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	cacheService := gocache.New(time.Hour*24, time.Hour)

	authMiddleware := middleware.NewAuthMiddleware(cacheService)
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository, cacheService)
	authHandler := handler.NewAuthHandler(authService)

	cloudinaryUtils, err := utils.NewCloudinaryUtils()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository, cloudinaryUtils)
	productHandler := handler.NewProductHandler(productService)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.ErrorMiddleware,
			authMiddleware.Middleware,
		),
	)

	auth.RegisterAuthServiceServer(server, authHandler)
	product.RegisterProductServiceServer(server, productHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(server)
		log.Println("reflection service registered")
	}

	wrappedGrpc := grpcweb.WrapServer(server,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"*",
			"Content-Type",
			"X-Grpc-Web",
			"X-User-Agent",
			"Authorization",
		},
		ExposedHeaders: []string{
			"Grpc-Status",
			"Grpc-Message",
			"Grpc-Status-Details-Bin",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: corsHandler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[gRPC-Web] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

			if wrappedGrpc.IsGrpcWebRequest(r) {
				wrappedGrpc.ServeHTTP(w, r)
				return
			}

			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}

			http.NotFound(w, r)
		})),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("gRPC-Web server listening on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("error serving HTTP server: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Panicf("error starting server: %v", err)
	}

	log.Println("gRPC server listening on :3000")
	if err := server.Serve(lis); err != nil {
		log.Panicf("error serving gRPC server: %v", err)
	}

}