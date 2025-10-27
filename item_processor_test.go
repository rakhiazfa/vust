package vust

import (
	"github.com/sirupsen/logrus"
)

type ExampleItemProcessor struct{}

func NewExampleItemProcessor() ItemProcessor {
	return &ExampleItemProcessor{}
}

func (p *ExampleItemProcessor) Process(ctx StepContext, log *logrus.Logger, item any) (any, error) {
	resource, ok := item.(*exampleResource)
	if !ok {
		return nil, nil
	}

	return resource, nil
}
