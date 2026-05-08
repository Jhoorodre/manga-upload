package progress

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

type ProgressTracker struct {
	Total int64
	Done  atomic.Int64
}

func (t *ProgressTracker) Increment() {
	t.Done.Add(1)
}

type Progress interface {
	Start(total int64, tracker *ProgressTracker)
	Finish(success bool)
}

func NewProgress(quiet bool) Progress {
	if quiet {
		return noopProgress{}
	}
	// Use stderr to not pollute stdout piping
	if isatty.IsTerminal(os.Stderr.Fd()) {
		return newBarProgress()
	}
	return newPlainProgress()
}

// -- Noop Progress --
type noopProgress struct{}

func (p noopProgress) Start(total int64, tracker *ProgressTracker) {}
func (p noopProgress) Finish(success bool)                         {}

// -- Plain Progress --
type plainProgress struct {
	tracker *ProgressTracker
	ticker  *time.Ticker
	done    chan struct{}
}

func newPlainProgress() *plainProgress {
	return &plainProgress{}
}

func (p *plainProgress) Start(total int64, tracker *ProgressTracker) {
	p.tracker = tracker
	p.ticker = time.NewTicker(2 * time.Second)
	p.done = make(chan struct{})

	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.ticker.C:
				fmt.Fprintf(os.Stderr, "Progresso: %d/%d\n", p.tracker.Done.Load(), p.tracker.Total)
			}
		}
	}()
}

func (p *plainProgress) Finish(success bool) {
	if p.ticker != nil {
		p.ticker.Stop()
		close(p.done)
	}
	fmt.Fprintf(os.Stderr, "Finalizado: %d/%d\n", p.tracker.Done.Load(), p.tracker.Total)
}

// -- Bar Progress --
type barProgress struct {
	tracker *ProgressTracker
	ticker  *time.Ticker
	done    chan struct{}
	prog    progress.Model
}

func newBarProgress() *barProgress {
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
	)
	return &barProgress{prog: prog}
}

func (p *barProgress) Start(total int64, tracker *ProgressTracker) {
	p.tracker = tracker
	p.ticker = time.NewTicker(100 * time.Millisecond)
	p.done = make(chan struct{})

	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.ticker.C:
				p.render()
			}
		}
	}()
}

func (p *barProgress) render() {
	done := float64(p.tracker.Done.Load())
	total := float64(p.tracker.Total)
	var percent float64
	if total > 0 {
		percent = done / total
	}

	width, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || width < 10 {
		width = 80
	}

	p.prog.Width = width - 30
	if p.prog.Width < 10 {
		p.prog.Width = 10
	}

	barStr := p.prog.ViewAs(percent)
	text := fmt.Sprintf(" %d/%d ", int64(done), int64(total))

	// \r clears the current line visually if followed by \x1b[K
	fmt.Fprintf(os.Stderr, "\r\033[K%s%s", barStr, text)
}

func (p *barProgress) Finish(success bool) {
	if p.ticker != nil {
		p.ticker.Stop()
		close(p.done)
	}
	p.render()
	fmt.Fprintln(os.Stderr)
}
