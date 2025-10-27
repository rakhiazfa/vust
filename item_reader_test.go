package vust

import (
	"io"

	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
)

type exampleResource struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
}

func generateExampleResources(count int) []*exampleResource {
	resources := make([]*exampleResource, count)

	for i := range count {
		resources[i] = &exampleResource{
			ID:       i + 1,
			Name:     faker.Name(),
			Username: faker.Username(),
			Email:    faker.Email(),
			Phone:    faker.Phonenumber(),
		}
	}

	return resources
}

var exampleResources = generateExampleResources(100000)

type ExampleItemReader struct {
	index int
}

func NewExampleItemReader() ItemReader {
	return &ExampleItemReader{
		index: 0,
	}
}

func (r *ExampleItemReader) Read(ctx StepContext, log *logrus.Logger) (any, error) {
	if r.index == len(exampleResources) {
		return nil, io.EOF
	}

	resource := exampleResources[r.index]
	r.index++

	return resource, nil
}
