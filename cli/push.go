package cli

import (
	"fmt"

	"github.com/blake-education/dogestry/image"
	"github.com/blake-education/dogestry/push"
	"github.com/blake-education/dogestry/remote"
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

	imageRoot, err := cli.WorkDir(imageIdentifier)
	if err != nil {
		return err
	}

	image := image.NewImage(imageIdentifier, imageRoot)

	remote, err := remote.NewRemote(remoteDef, cli.Config)
	if err != nil {
		return err
	}

	fmt.Println("remote", remote.Desc())

	if err := push.Push(image, remote); err != nil {
		return err
	}

	return nil
}
