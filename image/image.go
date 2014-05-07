package image

import (
	"github.com/blake-education/dogestry/remote"
	"github.com/blake-education/dogestry/repository"
)

type Image struct {
}

func (i *Image) NeedsPushToRemote(remote remote.Remote) (bool, error) {
}

func (i *Image) NeedsPullFromRemote(remote remote.Remote) (bool, error) {
}

func (i *Image) ImportRepository(repo *repository.Repository) error {
}
