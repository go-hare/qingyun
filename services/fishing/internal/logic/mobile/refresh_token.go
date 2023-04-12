package mobile

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	mobile_fishing "qingyun/services/fishing/proto/mobile"
)

func (m *MobileFishingService) RefreshToken(ctx context.Context, request *mobile_fishing.RefreshTokenRequest, response *mobile_fishing.RefreshTokenResponse) (err error) {
	response.Status = &mobile_fishing.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "RefreshToken",
	})

	if request.UserId < 0 {

	}
	fmt.Println(logger)
	return nil
}
