package gateway

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/hrest"
	"github.com/inter-hubly/pilot/server"
)

type Keeper interface {
	GetClientByPhoneNumberId(ctx context.Context, clientId string) (*valueobject.Client, error)
}

type keeperGateway struct {
	url string
}

func NewKeeper() *keeperGateway {

	var (
		keeperOnce sync.Once
		keeper     *keeperGateway
	)

	keeperOnce.Do(func() {
		keeper = &keeperGateway{
			url: server.GetGatewayHost().KeeperHost,
		}
	})
	return keeper
}

func (k *keeperGateway) GetClientByPhoneNumberId(ctx context.Context, phoneNumberId string) (*valueobject.Client, error) {
	hlog.Debug("keeperGateway.GetClientByPhoneNumberId", "PhoneNumberId", phoneNumberId)
	request := hrest.NewRequest(fmt.Sprintf("%s/api/client/%s/phone-number-id", k.url, phoneNumberId))
	err := request.CreateRequest(ctx, http.MethodGet)
	if err != nil {
		return nil, err
	}

	var voClient valueobject.Client

	if err = request.GetBody(&voClient); err != nil {
		return nil, err
	}
	return &voClient, nil
}
