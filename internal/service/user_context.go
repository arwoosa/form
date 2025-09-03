package service

import (
	"context"

	"github.com/arwoosa/vulpes/ezgrpc"
)

// UserInfo represents authenticated user information
type UserInfo struct {
	UserID     string
	UserEmail  string
	UserName   string
	UserAvatar string
	MerchantID string
}

// GetUserInfo extracts user information from context using Vulpes
func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	// Use Vulpes GetUser method to extract user info from context
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Vulpes User to our UserInfo struct
	userInfo := &UserInfo{
		UserID:     user.ID,
		UserEmail:  user.Email,
		UserName:   user.Name,
		UserAvatar: "", // Avatar field not available in ezgrpc.User
		MerchantID: user.Merchant,
	}

	return userInfo, nil
}

// ValidateUserAccess ensures user has access to the specified merchant
func ValidateUserAccess(userInfo *UserInfo, requiredMerchantID string) error {
	if userInfo.MerchantID != requiredMerchantID {
		return ErrUnauthorized
	}
	return nil
}
