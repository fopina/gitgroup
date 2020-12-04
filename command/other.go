package commands

import (
	"sync"

	"github.com/cheggaaa/pb/v3"
	flag "github.com/spf13/pflag"
)

// ProgressBar ...
type ProgressBar struct {
	bar *pb.ProgressBar
}

// Start ...
func (p *ProgressBar) Start(count int) {
	p.bar = pb.Full.Start(count)
}

// MultiThread ...
type MultiThread struct {
	wg      sync.WaitGroup
	threads int
	inputChannel,
	outputChannel chan interface{}
}

func (h *MultiThread) FlagSet() *flag.FlagSet {
	flags := flag.FlagSet{}
	flags.IntVarP(&h.threads, "threads", "t", 4, "number of working threads")
	return &flags
}

func (h *MultiThread) StartWorkers(fn func()) {
	h.wg.Add(h.threads)
	h.inputChannel = make(chan interface{}, h.threads)
	h.outputChannel = make(chan interface{}, h.threads)

	for i := 0; i < h.threads; i++ {
		go func() {
			defer h.wg.Done()
			fn()
		}()
	}
}

func (h *MultiThread) FeedWorkers(fn func()) {
	go func() {
		fn()
		close(h.inputChannel)
		h.wg.Wait()
		close(h.outputChannel)
	}()
}
