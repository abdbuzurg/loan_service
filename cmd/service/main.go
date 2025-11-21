package main

import (
	"loan_service/configs"
	"loan_service/internal/clients"
	"loan_service/internal/handler"
	"loan_service/internal/platform/database"
	messagebroker "loan_service/internal/platform/message_broker"
	loanpb "loan_service/internal/proto/loan"
	"loan_service/internal/repository"
	"loan_service/internal/usecase"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := configs.LoadConfig("../../configs")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	dbPool, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("DB connection failed: %s", err)
	}
	defer dbPool.Close()

	rabbitMQConn, err := messagebroker.NewRabbitMQConnection(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("RabbitMQ connection failed: %s", err)
	}
	defer rabbitMQConn.Close()

	queries := repository.New(dbPool)

	koinotAutoClient, err := clients.NewKoinotAutoClient(cfg.Clients.KoinotAuto)
	if err != nil {
		log.Fatalf("Failed to instantiate KoinotAuto client: %s", err)
	}

	asrLeasingClient, err := clients.NewAsrLeasingClient(cfg.Clients.AsrLeasing)
	if err != nil {
		log.Fatalf("Failed to instantiate ASR LEASING client: %s", err)
	}

	loanUC := usecase.New(queries, asrLeasingClient, koinotAutoClient)

	loanHandler := handler.New(loanUC)

	lis, err := net.Listen("tcp", cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	grpcServer := grpc.NewServer()

	loanpb.RegisterLoansServiceServer(grpcServer, loanHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
