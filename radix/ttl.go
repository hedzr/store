package radix

import (
	"context"
	"sync"
	"time"
)

type OnTTLRinging[T any] func(s *TTL[T], nd Node[T])

type TTL[T any] struct {
	treevec []*trieS[T]
	mu      sync.RWMutex
	cancel  context.CancelFunc // exit signal
	done    <-chan struct{}
	adder   chan ttljobS[T]

	// in the future, we might try timing-wheel way.

	weeks   *wheelS // by weeks, years, ...
	years   *wheelS // by years, months, days, weeks
	days    *wheelS // by days, hours, minutes, seconds
	seconds *wheelS // by seconds, ms, us, ns
}

type wheelS struct {
	wheel    map[int]*wheelS
	mu       sync.RWMutex
	onAction wheelAction
}

type wheelAction func(ctx context.Context, w *wheelS)

type ttljobS[T any] struct {
	node     *nodeS[T]
	duration time.Duration
	action   OnTTLRinging[T]
}

func newttls[T any](t *trieS[T]) *TTL[T] {
	ctx, cancel := context.WithCancel(context.Background())
	s := &TTL[T]{
		treevec: []*trieS[T]{t},
		cancel:  cancel,
		done:    ctx.Done(),
		adder:   make(chan ttljobS[T]),
	}
	go s.run()
	return s
}

func (s *TTL[T]) dupS() *TTL[T] {
	ctx, cancel := context.WithCancel(context.Background())
	n := &TTL[T]{
		treevec: []*trieS[T]{s.treevec[0]},
		cancel:  cancel,
		done:    ctx.Done(),
		adder:   make(chan ttljobS[T]),
	}
	go n.run()
	return n
}

func (s *TTL[T]) Close() {
	s.cancel()
}

func (s *TTL[T]) Tree() Trie[T] { return s.treevec[0] }

func (s *TTL[T]) Add(nd *nodeS[T], duration time.Duration, action OnTTLRinging[T]) {
	if nd.rw == nil {
		nd.rw = &sync.RWMutex{}
	}
	s.adder <- ttljobS[T]{node: nd, duration: duration, action: action}
}

func (s *TTL[T]) run() {
	adder := func(job ttljobS[T]) {
		if job.duration < 200*time.Nanosecond {
			if job.duration == 0 {
				// reset a job (might be planned in the future)
			}
			return
		}

		// add a new job
		timer := time.NewTimer(job.duration)
		go func(timer *time.Timer, job ttljobS[T]) {
			defer timer.Stop()
			// defer func() { job.node.rw = nil }()
			for {
				select {
				case <-timer.C:
					if job.action != nil {
						job.action(s, job.node)
					}
					if job.node.isBranch() {
						s.treevec[0].Remove(job.node.pathS)
					} else {
						job.node.SetEmpty()
					}
					return
				case <-s.done:
					return
				}
			}
		}(timer, job)
	}

	for {
		select {
		case <-s.done:
			return
		case job := <-s.adder:
			adder(job)
		}
	}
}
