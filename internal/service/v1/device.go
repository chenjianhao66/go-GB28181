package v1

import "github.com/chenjianhao66/go-GB28181/internal/store"

type DeviceSrv interface {
}

type deviceService struct {
	store store.Factory
}

func newDeviceSrv(store store.Factory) DeviceSrv {
	return &deviceService{
		store: store,
	}
}
