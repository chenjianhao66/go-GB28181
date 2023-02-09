package storage

import "github.com/chenjianhao66/go-GB28181/internal/model"

type MediaStorage interface {
	GetMediaByID(id string) (model.MediaDetail, error)
	Save(config model.MediaDetail) error
}
