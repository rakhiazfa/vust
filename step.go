package vust

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type Step interface {
	WithContext(StepContext) Step
	SetLogger(*logrus.Logger) Step
	SetJobName(name string) Step
	Run() error
}

type step struct {
	config  StepConfig
	ctx     StepContext
	log     *logrus.Logger
	jobName string
}

type StepConfig struct {
	Name          string
	Reader        ItemReader
	Processor     ItemProcessor
	Writer        ItemWriter
	ChunkSize     int
	ErrorListener ErrorListener
}

func NewStep(config ...StepConfig) Step {
	s := &step{
		config: StepConfig{},
		ctx:    NewStepContext(NewJobContext(context.Background())),
		log:    logrus.New(),
	}

	if len(config) > 0 {
		s.config = config[0]
	}

	if s.config.ErrorListener == nil {
		s.config.ErrorListener = NewDefaultErrorListener()
	}

	return s
}

func (s *step) WithContext(ctx StepContext) Step {
	newStep := *s
	newStep.ctx = ctx

	return &newStep
}

func (s *step) SetLogger(log *logrus.Logger) Step {
	s.log = logrus.New()

	s.log.SetOutput(log.Out)
	s.log.SetLevel(log.Level)
	s.log.SetReportCaller(log.ReportCaller)
	s.log.SetFormatter(log.Formatter)

	s.log.ExitFunc = log.ExitFunc
	s.log.BufferPool = log.BufferPool

	return s
}

func (s *step) SetJobName(name string) Step {
	s.jobName = name
	return s
}

func (s *step) Run() error {
	if err := s.validate(); err != nil {
		return err
	}

	s.log.AddHook(&StepPrefixHook{Name: s.config.Name})
	s.log.AddHook(&JobPrefixHook{Name: s.jobName})
	s.log.AddHook(&ModulePrefixHook{})

	batchCh := make(chan *Batch, 5)

	var stepWg sync.WaitGroup

	stepWg.Go(func() {
		itemCh := s.readItems()
		processedItemCh := s.processItems(itemCh)

		batchId := 1
		items := make([]any, 0, s.config.ChunkSize)

		for processedItem := range processedItemCh {
			items = append(items, processedItem)

			if len(items) == s.config.ChunkSize {
				batch := &Batch{ID: batchId, Items: items}
				batchId++

				batchCh <- batch

				items = items[:0]
			}
		}

		if len(items) > 0 {
			batch := &Batch{ID: batchId, Items: items}

			batchCh <- batch
		}

		close(batchCh)
	})
	stepWg.Go(func() {
		for batch := range batchCh {
			s.writeItems(batch)
		}
	})

	stepWg.Wait()

	return nil
}

func (s *step) readItems() <-chan any {
	itemCh := make(chan any)

	go func() {
		for {
			item, err := s.config.Reader.Read(s.ctx, s.log)
			if err != nil {
				if err == io.EOF {
					break
				}

				s.handleErrorOnReading(err)
				continue
			}

			if item == nil {
				continue
			}

			itemCh <- item
		}

		close(itemCh)
	}()

	return itemCh
}

func (s *step) processItems(itemCh <-chan any) <-chan any {
	processedItemCh := make(chan any)

	go func() {
		for item := range itemCh {
			processed, err := s.config.Processor.Process(s.ctx, s.log, item)
			if err != nil {
				s.handleErrorOnProcessing(item, err)
				continue
			}

			if processed == nil {
				continue
			}

			processedItemCh <- processed
		}

		close(processedItemCh)
	}()

	return processedItemCh
}

func (s *step) writeItems(batch *Batch) {
	if err := s.config.Writer.Write(s.ctx, s.log, batch); err != nil {
		s.handleErrorOnWriting(batch, err)
		return
	}
}

func (s *step) validate() error {
	if strings.TrimSpace(s.config.Name) == "" {
		return errors.New("step name cannot be empty")
	}
	if s.config.Reader == nil {
		return errors.New("reader cannot be empty")
	}
	if s.config.Processor == nil {
		return errors.New("processor cannot be empty")
	}
	if s.config.Writer == nil {
		return errors.New("writer cannot be empty")
	}
	if s.config.ChunkSize <= 0 {
		return errors.New("chunk size cannot be less than equal to zero")
	}

	return nil
}

func (s *step) handleErrorOnReading(err error) {
	if s.config.ErrorListener != nil {
		s.config.ErrorListener.OnRead(s.ctx, s.log, err)
	}
}

func (s *step) handleErrorOnProcessing(item any, err error) {
	if s.config.ErrorListener != nil {
		s.config.ErrorListener.OnProcess(s.ctx, s.log, item, err)
	}
}

func (s *step) handleErrorOnWriting(batch *Batch, err error) {
	if s.config.ErrorListener != nil {
		s.config.ErrorListener.OnWrite(s.ctx, s.log, batch, err)
	}
}
