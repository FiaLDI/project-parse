package output

import (
	"context"
	"errors"

	"github.com/FiaLDI/project-parse/internal/domain"
)

// ErrNotImplemented indicates a renderer is reserved but not yet available.
var ErrNotImplemented = errors.New("renderer not implemented")

// PDFRenderer is a placeholder for future PDF export support.
type PDFRenderer struct{}

func NewPDF() *PDFRenderer { return &PDFRenderer{} }

func (r *PDFRenderer) Format() string { return "pdf" }

func (r *PDFRenderer) Render(_ context.Context, _ domain.RenderDocument) ([]byte, error) {
	return nil, ErrNotImplemented
}
