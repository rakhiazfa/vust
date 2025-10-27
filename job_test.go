package vust

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestJob(t *testing.T) {
	job := New(JobConfig{
		Name: "Customer Import Job",
	})

	customerImportStep := NewStep(StepConfig{
		Name:      "Customer Import Step",
		Reader:    NewExampleItemReader(),
		Processor: NewExampleItemProcessor(),
		Writer:    NewExampleItemWriter(),
		ChunkSize: 1000,
	})
	transactionImportStep := NewStep(StepConfig{
		Name:      "Transaction Import Step",
		Reader:    NewExampleItemReader(),
		Processor: NewExampleItemProcessor(),
		Writer:    NewExampleItemWriter(),
		ChunkSize: 1000,
	})

	job.AddStep(customerImportStep)
	job.AddStep(transactionImportStep)

	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	ctx := NewJobContext(context.Background())
	ctx.Set("groupId", "GRP-27/10/25-001")

	ctx, cancel := ctx.WithTimeout(1 * time.Second)
	defer cancel()

	err := job.WithContext(ctx).SetLogger(log).Run()

	assert.NoError(t, err)
}
