package service

import (
	"context"

	"github.com/arwoosa/form-service/pkg/vulpes/relation"
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
	user, err := relation.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Vulpes User to our UserInfo struct
	userInfo := &UserInfo{
		UserID:     user.UserID,
		UserEmail:  user.Email,
		UserName:   user.Name,
		UserAvatar: user.Avatar,
		MerchantID: user.MerchantID,
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