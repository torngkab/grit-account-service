package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	pb "github.com/torngkab/grit-account-service/account"
	"github.com/torngkab/grit-account-service/config"
	"github.com/torngkab/grit-account-service/database"
	"github.com/torngkab/grit-account-service/model"
	"github.com/torngkab/grit-account-service/utils"
	"gorm.io/gorm"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAccountServer
	gorm *gorm.DB
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// declare user
	user := model.User{
		Id:        uuid.New(),
		Username:  in.GetUsername(),
		Password:  utils.HashPassword(in.GetPassword()),
		Status:    model.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.gorm.Transaction(func(tx *gorm.DB) error {
		// Create user
		if err := tx.Create(&user).Error; err != nil {
			log.Printf("failed to create user: %v", err)
			if strings.Contains(err.Error(), "duplicate key") {
				return errors.New("user already exists")
			}
			return err
		}

		// Create accounts
		accountId := uuid.New()
		accounts := []model.Account{
			{
				Id:          accountId,
				UserId:      user.Id,
				AccountType: model.AccountTypeMain,
				Status:      model.AccountStatusActive,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		if err := tx.Create(&accounts).Error; err != nil {
			log.Printf("failed to create accounts: %v", err)
			return err
		}

		// Create balances
		balances := []model.Balance{
			{
				AccountId:       accountId,
				Balance:         0,
				LatestUpdatedAt: time.Now(),
			},
		}
		if err := tx.Create(&balances).Error; err != nil {
			log.Printf("failed to create balances: %v", err)
			return err
		}

		return nil
	})
	if err != nil {
		return &pb.CreateUserResponse{Message: err.Error()}, nil
	}

	// TODO: Send referral code to referral service...

	return &pb.CreateUserResponse{Message: "user created successfully"}, nil
}

func (s *server) GetAccountsByUserId(ctx context.Context, in *pb.GetAccountsByUserIdRequest) (*pb.GetAccountsByUserIdResponse, error) {
	// declare accounts
	accounts := []model.Account{}

	// get accounts
	err := s.gorm.Where("user_id = ?", in.GetUserId()).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	// check if accounts is empty
	if len(accounts) == 0 {
		return &pb.GetAccountsByUserIdResponse{Accounts: []*pb.AccountModel{}}, nil
	}

	// convert accounts to protobuf
	accountsPb := make([]*pb.AccountModel, len(accounts))
	for i, account := range accounts {
		accountsPb[i] = &pb.AccountModel{
			Id:          account.Id.String(),
			UserId:      account.UserId.String(),
			AccountType: string(account.AccountType),
			Status:      string(account.Status),
			CreatedAt:   account.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   account.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &pb.GetAccountsByUserIdResponse{Accounts: accountsPb}, nil
}

func (s *server) GetAccountBalance(ctx context.Context, in *pb.GetAccountBalanceRequest) (*pb.GetAccountBalanceResponse, error) {
	// declare balance
	balance := model.Balance{}

	// get balance
	err := s.gorm.Where("account_id = ?", in.GetAccountId()).First(&balance).Error
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountBalanceResponse{
		Balance:         float32(balance.Balance),
		LatestUpdatedAt: balance.LatestUpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *server) SetAccountBalance(ctx context.Context, in *pb.SetAccountBalanceRequest) (*pb.SetAccountBalanceResponse, error) {
	// declare balance
	balance := model.Balance{}

	// get balance
	err := s.gorm.Where("account_id = ?", in.GetAccountId()).First(&balance).Error
	if err != nil {
		log.Printf("failed to get balance: %v", err)
		return nil, err
	}

	// declare new balance
	newBalance := model.Balance{
		Balance:         balance.Balance + float64(in.GetBalance()),
		LatestUpdatedAt: time.Now(),
	}

	// set balance
	err = s.gorm.Model(&model.Balance{}).
		Where("account_id = ?", in.GetAccountId()).
		Updates(map[string]interface{}{
			"balance":           newBalance.Balance,
			"latest_updated_at": newBalance.LatestUpdatedAt,
		}).
		Error
	if err != nil {
		log.Printf("failed to set balance: %v", err)
		return nil, err
	}

	return &pb.SetAccountBalanceResponse{Balance: float32(newBalance.Balance), LatestUpdatedAt: newBalance.LatestUpdatedAt.Format(time.RFC3339)}, nil
}

func (s *server) ValidateAccountBalance(ctx context.Context, in *pb.ValidateAccountBalanceRequest) (*pb.ValidateAccountBalanceResponse, error) {
	// declare balance
	balance := model.Balance{}

	// get balance
	err := s.gorm.Where("account_id = ?", in.GetAccountId()).
		First(&balance).Error
	if err != nil {
		return nil, err
	}

	return &pb.ValidateAccountBalanceResponse{IsValid: balance.Balance >= float64(in.GetAmount())}, nil
}

func main() {
	// Get config
	config := config.C("")

	// Connect to PostgreSQL
	gormDB := database.NewGormDB(config.Database.PostgresURL)

	if config.IsLocalEnv() {
		// Auto migrate models
		gormDB.AutoMigrate(&model.User{}, &model.Account{}, &model.Balance{})

		gormDB.Transaction(func(tx *gorm.DB) error {
			// Seed data
			userId := uuid.New()
			accountIds := []uuid.UUID{
				uuid.New(), // Referral
				uuid.New(), // Main
				uuid.New(), // Disbursement
				uuid.New(), // PSP
			}

			// Create user
			tx.Create(&model.User{
				Id:        userId,
				Username:  "admin",
				Password:  utils.HashPassword("admin"),
				Status:    model.UserStatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})

			// Create accounts
			tx.Create([]model.Account{
				{
					Id:          accountIds[0],
					UserId:      userId,
					AccountType: model.AccountTypeReferral,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Id:          accountIds[1],
					UserId:      userId,
					AccountType: model.AccountTypeMain,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Id:          accountIds[2],
					UserId:      userId,
					AccountType: model.AccountTypeDisbursement,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Id:          accountIds[3],
					UserId:      userId,
					AccountType: model.AccountTypePSP,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			})

			// Create balances
			tx.Create([]model.Balance{
				{
					AccountId:       accountIds[0],
					Balance:         10000,
					LatestUpdatedAt: time.Now(),
				},
				{
					AccountId:       accountIds[1],
					Balance:         0,
					LatestUpdatedAt: time.Now(),
				},
				{
					AccountId:       accountIds[2],
					Balance:         0,
					LatestUpdatedAt: time.Now(),
				},
				{
					AccountId:       accountIds[3],
					Balance:         0,
					LatestUpdatedAt: time.Now(),
				},
			})

			return nil
		})
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterAccountServer(s, &server{
		gorm: gormDB,
	})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
