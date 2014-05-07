package cli

import (
	"fmt"

	"github.com/blake-education/dogestry/actions"
	"github.com/blake-education/dogestry/image"
	"github.com/blake-education/dogestry/remote"
	"github.com/blake-education/dogestry/repository"
)

func (cli *DogestryCli) CmdPush(args ...string) error {
	cmd := cli.Subcmd("push", "REMOTE IMAGE[:TAG]", "push IMAGE to the REMOTE. TAG defaults to 'latest'")
	if err := cmd.Parse(args); err != nil {
		return nil
	}

	if len(cmd.Args()) < 2 {
		return fmt.Errorf("Error: IMAGE and REMOTE not specified")
	}

	remoteDef := cmd.Arg(0)
	imageIdentifier := cmd.Arg(1)

	repoRoot, err := cli.WorkDir(imageIdentifier)
	if err != nil {
		return err
	}

	image := image.NewImage(cli.client, imageIdentifier)
	repo := repository.NewRepo(repoRoot)

	remote, err := remote.NewRemote(remoteDef, cli.Config)
	if err != nil {
		return err
	}

	fmt.Println("remote", remote.Desc())

	if err := actions.Push(image, repository, remote); err != nil {
		return err
	}

	return nil
}
