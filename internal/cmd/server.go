package cmd

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	pbComment "github.com/kodinggo/comment-service-gp1/pb/comment"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/config"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/db"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/delivery/grpcsvc"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/delivery/httpsvc"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/repository"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/usecase"
	pb "github.com/kodinggo/rest-api-service-golang-private-1/pb/story"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverCMD = &cobra.Command{
	Use:   "serve",
	Short: "serve is a command to run the service",
	Run: func(cmd *cobra.Command, args []string) {
		// Init DB connection
		dbConn := db.InitMySQLConn()

		sigCh := make(chan os.Signal, 1)
		errCh := make(chan error, 1)
		quitCh := make(chan bool, 1)
		signal.Notify(sigCh, os.Interrupt)

		go func() {
			for i := 0; i < 2; i++ {
				select {
				case <-sigCh:
					log.Println("shutdown due to interrupt signal")
					quitCh <- true
				case err := <-errCh:
					log.Printf("failed when running server, error: %v", err)
					quitCh <- true
				}
			}
		}()

		// Init gRPC clients
		grpcCommentClient := newgRPCCommentClient()

		// Init Repository
		storyRepo := repository.NewStoryRepository(dbConn)
		userRepo := repository.NewUserRepository(dbConn)

		// Init Usecase
		storyUsecase := usecase.NewStoryUsecase(storyRepo, userRepo, grpcCommentClient)
		authUsecase := usecase.NewAuthUsecase(userRepo)

		go func() {
			// Run HTTP server
			e := echo.New()

			// It's used to verify
			e.GET("/ping", func(c echo.Context) error {
				return c.String(http.StatusOK, "pong!")
			})

			// Init HTTP Handler
			storyHandler := httpsvc.NewStoryHandler(storyUsecase)
			authHandler := httpsvc.NewAuthHandler(authUsecase)

			storyHandler.RegisterRoutes(e)
			authHandler.RegisterRoutes(e)

			errCh <- e.Start(config.Port())
		}()

		go func() {
			// Run grpc server
			grpcServer := grpc.NewServer()

			storyService := grpcsvc.NewStoryService(storyUsecase)
			userService := grpcsvc.NewUserService()

			pb.RegisterStoryServiceServer(grpcServer, storyService)
			pb.RegisterUserServiceServer(grpcServer, userService)

			lis, err := net.Listen("tcp", config.GRPCPort())
			if err != nil {
				errCh <- err
				return
			}

			log.Printf("grpc server is running with port: %s", config.GRPCPort())

			err = grpcServer.Serve(lis)
			if err != nil {
				errCh <- err
			}
		}()

		<-quitCh
		log.Printf("exiting")
	},
}

func init() {
	rootCMD.AddCommand(serverCMD)
}

func newgRPCCommentClient() pbComment.CommentServiceClient {
	// connect to grpc server without credentials
	conn, err := grpc.NewClient("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("failed to open connection grpc server, error %v", err)
	}
	// init grpc client as package dependency from grpc-server repository
	return pbComment.NewCommentServiceClient(conn)
}
