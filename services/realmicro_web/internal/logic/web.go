package logic

import (
	"context"
	"qingyun/services/realmicro_web/internal/svc"
	"qingyun/services/realmicro_web/proto"
)

type MicroWeb struct {
	SvcContext *svc.ServiceContext
}

func (m *MicroWeb) ProjectAlias(ctx context.Context, request *micro_web.ProjectAliasRequest, response *micro_web.ProjectAliasResponse) error {
	response.Status = &micro_web.Status{}
	return nil
}

func NewMicroWeb(ctx *svc.ServiceContext) (microWeb *MicroWeb) {
	microWeb = &MicroWeb{
		SvcContext: ctx,
	}
	return
}
