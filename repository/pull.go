package repository

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"

	"github.com/blake-education/dogestry/remote"
)

func (r *Repository) PullImageFromRemote(image *image.Image, remote remote.Remote) error {
}

func prepareRepositories(image, imageRoot string, r remote.Remote) error {
	repoName, repoTag := remote.NormaliseImageName(image)

	id, err := r.ParseTag(repoName, repoTag)
	if err != nil {
		return err
	} else if id == "" {
		return nil
	}

	reposPath := filepath.Join(imageRoot, "repositories")
	reposFile, err := os.Create(reposPath)
	if err != nil {
		return err
	}
	defer reposFile.Close()

	repositories := map[string]Repository{}
	repositories[repoName] = Repository{}
	repositories[repoName][repoTag] = string(id)

	return json.NewEncoder(reposFile).Encode(&repositories)
}
