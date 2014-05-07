package image

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/fsouza/go-dockerclient/utils"
)

// Stream the tarball from docker and translate it into the portable repo format
// Note that its easier to handle as a stream on the way out.
func (i *Image) ExtractImageLayers(client docker.Client) error {
	reader, writer := io.Pipe()
	defer writer.Close()
	defer reader.Close()

	tarball := tar.NewReader(reader)

	errch := make(chan error)

	go func() {
		// consume the tar
		for {
			header, err := tarball.Next()
			if err == io.EOF {
				// end of tar archive
				break
			}
			if err != nil {
				errch <- err
				return
			}

			if err := p.processTarEntry(header, tarball); err != nil {
				errch <- err
				return
			}
		}

		// donno... read a bit more?
		if _, err := ioutil.ReadAll(reader); err != nil {
			errch <- err
			return
		}

		errch <- nil
	}()

	if err := client.GetImageTarball(image, writer); err != nil {
		// this should stop the tar reader
		writer.Close()
		<-errch
		return err
	}

	writer.Close()

	// wait for the tar reader
	if err := <-errch; err != nil {
		return err
	}

	return nil
}

func (i *Image) processTarEntry(header *tar.Header, tarball io.Reader) error {
	// only handle files (directories are implicit)
	if header.Typeflag == tar.TypeReg {
		fmt.Printf("  tar: processing %s\n", header.Name)

		// special case - repositories file
		if filepath.Base(header.Name) == "repositories" {
			if err := writeRepositories(i.localRoot, tarball); err != nil {
				return err
			}

		} else {
			barename := strings.TrimPrefix(header.Name, "./")

			dest := filepath.Join(i.localRoot, "images", barename)
			if err := os.MkdirAll(filepath.Dir(dest), os.ModeDir|0700); err != nil {
				return err
			}

			destFile, err := os.Create(dest)
			if err != nil {
				return err
			}
			defer destFile.Close()

			if wrote, err := io.Copy(destFile, tarball); err != nil {
				return err
			} else {
				fmt.Printf("  tar: wrote %s\n", utils.HumanSize(wrote))
			}
			destFile.Close()
		}
	}

	return nil
}

type Repository map[string]string

func writeRepositories(root string, tarball io.Reader) error {
	destRoot := filepath.Join(root, "repositories")

	repositories := map[string]Repository{}
	if err := json.NewDecoder(tarball).Decode(&repositories); err != nil {
		return err
	}

	for repoName, repo := range repositories {
		for tag, id := range repo {
			dest := filepath.Join(destRoot, repoName, tag)

			if err := os.MkdirAll(filepath.Dir(dest), os.ModeDir|0700); err != nil {
				return err
			}

			if err := ioutil.WriteFile(dest, []byte(id), 0600); err != nil {
				return err
			}
		}
	}

	return nil
}
