package referralservice

import (
	"context"

	"github.com/torngkab/grit-account-service/config"
	"google.golang.org/grpc"

	pb "github.com/torngkab/grit-referral-service/referral"
)

type ReferralServiceAdapter interface {
	CreateReferral(ctx context.Context, userId string) (*pb.CreateReferralResponse, error)
	UseReferral(ctx context.Context, referralCode string, userId string) (*pb.UseReferralResponse, error)
}

type referralServiceAdapter struct {
	config config.Config
}

func NewReferralServiceAdapter(config config.Config) ReferralServiceAdapter {
	return &referralServiceAdapter{config: config}
}

func (a *referralServiceAdapter) CreateReferral(ctx context.Context, userId string) (*pb.CreateReferralResponse, error) {
	conn, err := grpc.Dial(a.config.Service.ReferralService, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewReferralClient(conn)
	resp, err := client.CreateReferral(ctx, &pb.CreateReferralRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *referralServiceAdapter) UseReferral(ctx context.Context, referralCode string, userId string) (*pb.UseReferralResponse, error) {
	conn, err := grpc.Dial(a.config.Service.ReferralService, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewReferralClient(conn)
	resp, err := client.UseReferral(ctx, &pb.UseReferralRequest{
		ReferralCode: referralCode,
		UserId:       userId,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
