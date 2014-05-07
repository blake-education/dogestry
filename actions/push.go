package actions

import (
	"fmt"

	"github.com/blake-education/dogestry/image"
	"github.com/blake-education/dogestry/remote"
	"github.com/blake-education/dogestry/repository"
)

func Push(image *image.Image, remote remote.Remote, repo *repository.Repository) error {
	needsPush, err := image.NeedsPushToRemote(remote)
	if err != nil {
		return err
	}
	if !needsPush {
		fmt.Println("push: remote already up to date")
		return nil
	}

	fmt.Println("push: preparing image")
	if err := image.ExtractImageLayersIntoRepository(repo); err != nil {
		return err
	}

	fmt.Println("push: pushing image to remote")
	if err := repo.PushToRemote(remote); err != nil {
		return err
	}
}
