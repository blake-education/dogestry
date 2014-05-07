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



	fmt.Println("resolving image id")
	id, err := r.ResolveImageNameToId(image)
	if err != nil {
		return err
	}

	fmt.Printf("image '%s' resolved on remote id '%s'\n", image, id.Short())

	fmt.Println("preparing images")
	if err := cli.preparePullImage(id, imageRoot, r); err != nil {
		return err
	}

	fmt.Println("preparing repositories file")
	if err := prepareRepositories(image, imageRoot, r); err != nil {
		return err
	}

	fmt.Println("sending tar to docker")
	if err := cli.sendTar(imageRoot); err != nil {
		return err
	}
