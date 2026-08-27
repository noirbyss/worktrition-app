package grpcclient

import (
	"ai-service/internal/domain"
	"context"
	"fmt"

	"github.com/noirbyss/worktrition-app/gen/user-service"
	"go.uber.org/zap"
)

type UserClient struct {
	client user.UserServiceClient
	logger *zap.SugaredLogger
}

func NewUserClient (client user.UserServiceClient, logger *zap.SugaredLogger) *UserClient {
	return &UserClient{
		client: client,
		logger: logger,
	}
} 

func (uc *UserClient) GetProfile(ctx context.Context, userId string) (*domain.UserProfile, error) {
	resp, err := uc.client.GetProfile(ctx, &user.GetProfileRequest{UserId: userId})
	if err != nil{
		uc.logger.Errorf("failed to getProfile from user-service for user_id=%s: %v",userId ,err)
		return nil, fmt.Errorf("user client getProfile %w",err)
	}
	profile := mapProtoToUserProfile(resp)
	
	return profile,nil
}

