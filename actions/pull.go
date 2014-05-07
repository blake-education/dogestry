package actions

import (
	"fmt"

	"github.com/blake-education/dogestry/image"
	"github.com/blake-education/dogestry/remote"
	"github.com/blake-education/dogestry/repository"
)

func Pull(image *image.Image, remote remote.Remote, repo *repository.Repository) error {
	needsPull, err := image.NeedsPullFromRemote(remote)
	if err != nil {
		return err
	}
	if !needsPull {
		fmt.Println("pull: local docker already up to date")
		return nil
	}

	err := repo.PullImageFromRemote(image, remote)
	if err != nil {
		return err
	}

	err := image.ImportRepository(repo)
	if err != nil {
		return err
	}

	return nil
}
