package mobile

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	mobile_fishing_user "qingyun/services/fishing_user/proto/mobile"
)

func (m *MobileFishingUserService) RefreshToken(ctx context.Context, request *mobile_fishing_user.RefreshTokenRequest, response *mobile_fishing_user.RefreshTokenResponse) (err error) {
	response.Status = &mobile_fishing_user.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "RefreshToken",
	})

	if request.UserId < 0 {

	}
	fmt.Println(logger)
	return nil
}
