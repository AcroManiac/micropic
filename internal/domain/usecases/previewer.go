package usecases

import (
	"net/http"

	"github.com/AcroManiac/micropic/internal/domain/entities"
	"github.com/AcroManiac/micropic/internal/domain/interfaces"
)

// Previewer object for image operations
type Previewer struct {
	loader    interfaces.ImageLoader
	processor interfaces.ImageProcessor
	senders   []interfaces.Sender
}

// NewPreviewer constructor
func NewPreviewer(
	loader interfaces.ImageLoader,
	processor interfaces.ImageProcessor,
	senders []interfaces.Sender,
) *Previewer {
	return &Previewer{
		loader:    loader,
		processor: processor,
		senders:   senders,
	}
}

func (p *Previewer) send(response *entities.Response, objs ...interface{}) {
	for _, s := range p.senders {
		s.Send(response, objs)
	}
}

// Process proxy request to preview response
func (p *Previewer) Process(request *entities.Request, objs ...interface{}) {
	// Load image from external source
	srcImage, status := p.loader.Load(request)
	if status.Code != http.StatusOK {
		response := entities.NewResponse(srcImage, "", *status)
		p.send(response, objs)
		return
	}

	// Make preview from source image
	response := p.processor.Process(srcImage, request)

	// Send preview response to consumers
	p.send(response, objs)
}