package cli

import (
	"fmt"

	"github.com/blake-education/dogestry/actions"
	"github.com/blake-education/dogestry/remote"
)

func (cli *DogestryCli) CmdPull(args ...string) error {
	cmd := cli.Subcmd("pull", "REMOTE IMAGE[:TAG]", "pull IMAGE from the REMOTE and load it into docker. TAG defaults to 'latest'")
	if err := cmd.Parse(args); err != nil {
		return nil
	}

	if len(cmd.Args()) < 2 {
		return fmt.Errorf("Error: REMOTE and IMAGE not specified")
	}

	remoteDef := cmd.Arg(0)
	image := cmd.Arg(1)

	workRoot, err := cli.WorkDir(image)
	if err != nil {
		return err
	}

	repository := repository.NewRepo(workRoot)

	r, err := remote.NewRemote(remoteDef, cli.Config)
	if err != nil {
		return err
	}

	fmt.Println("remote", r.Desc())

	image := image.NewImage(image, cli.Client)

	return actions.Pull(image, remote, repo)
}
