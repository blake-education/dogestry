package push

import (
	"fmt"

	"github.com/blake-education/dogestry/image"
	"github.com/blake-education/dogestry/remote"
)

type Push struct {
	image     string
	imageRoot string
	remote    remote.Remote
}

func Push(image *image.Image, remote remote.Remote) error {
	p := &Push{
		image:     image,
		imageRoot: imageRoot,
		remote:    remote,
	}

	fmt.Println("preparing image")
	if err := image.ExtractImageLayers(p.client); err != nil {
		return err
	}

	fmt.Println("pushing image to remote")
	if err := p.push(); err != nil {
		return err
	}
}

func (p *Push) push() error {
	p.remote.Push(p.image, p.imageRoot)
}
