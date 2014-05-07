package image

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/blake-education/dogestry/remote"
	"github.com/blake-education/dogestry/repository"
	"github.com/fsouza/go-dockerclient"
)

type Image struct {
}

func (i *Image) NeedsPushToRemote(remote remote.Remote) (bool, error) {
}

func (i *Image) NeedsPullFromRemote(remote remote.Remote) (bool, error) {
}

func (i *Image) ImportRepository(repo *repository.Repository) error {
}

func (cli *DogestryCli) preparePullImage(fromId remote.ID, imageRoot string, r remote.Remote) error {
	toDownload := make([]remote.ID, 0)

	// TODO flatten this list, then iterate and pull each required file
	// TODO parallelize
	err := r.WalkImages(fromId, func(id remote.ID, image docker.Image, err error) error {
		fmt.Printf("examining id '%s' on remote\n", id.Short())
		if err != nil {
			fmt.Println("err", err)
			return err
		}

		_, err = cli.client.InspectImage(string(id))
		if err == docker.ErrNoSuchImage {
			toDownload = append(toDownload, id)
			return nil
		} else if err != nil {
			return err
		} else {
			fmt.Printf("docker already has id '%s', stopping\n", id.Short())
			return remote.BreakWalk
		}
	})

	if err != nil {
		return err
	}

	for _, id := range toDownload {
		if err := cli.pullImage(id, filepath.Join(imageRoot, string(id)), r); err != nil {
			return err
		}
	}

	return nil
}

func (cli *DogestryCli) pullImage(id remote.ID, dst string, r remote.Remote) error {
	fmt.Printf("pulling image id '%s'\n", id.Short())

	// XXX fix image name rewrite
	err := r.PullImageId(id, dst)
	if err != nil {
		return err
	}
	return cli.processPulled(id, dst)
}

// no-op for now
func (cli *DogestryCli) processPulled(id remote.ID, dst string) error {
	return nil
}

// stream the tarball into docker
// its easier here to use tar command, but it'd be neater to mirror Push's approach
func (cli *DogestryCli) sendTar(imageRoot string) error {
	notExist, err := dirNotExistOrEmpty(imageRoot)

	if err != nil {
		return err
	}
	if notExist {
		fmt.Println("no images to send to docker")
		return nil
	}

	// DEBUG - write out a tar to see what's there!
	// exec.Command("/bin/tar", "cvf", "/tmp/d.tar", "-C", imageRoot, ".").Run()

	cmd := exec.Command("/bin/tar", "cvf", "-", "-C", imageRoot, ".")
	cmd.Dir = imageRoot
	defer cmd.Wait()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Println("kicking off post")
	return cli.client.PostImageTarball(stdout)
}

func dirNotExistOrEmpty(path string) (bool, error) {
	imagesDir, err := os.Open(path)
	if err != nil {
		// no images
		if os.IsNotExist(err) {
			return true, nil
		} else {
			return false, err
		}
	}
	defer imagesDir.Close()

	names, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	if len(names) <= 1 {
		return true, nil
	}

	return false, nil
}
