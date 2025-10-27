package vust

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Job interface {
	WithContext(JobContext) Job
	SetLogger(*logrus.Logger) Job
	AddStep(Step) Job
	Run() error
}

type job struct {
	config JobConfig
	ctx    JobContext
	log    *logrus.Logger
	steps  []Step
}

type JobConfig struct {
	Name        string
	JobListener JobListener
}

func New(config ...JobConfig) Job {
	j := &job{
		config: JobConfig{},
		ctx:    NewJobContext(context.Background()),
		log:    logrus.New(),
	}

	if len(config) > 0 {
		j.config = config[0]
	}

	if j.config.JobListener == nil {
		j.config.JobListener = NewDefaultJobListener()
	}

	return j
}

func (j *job) WithContext(ctx JobContext) Job {
	newJob := *j
	newJob.ctx = ctx

	return &newJob
}

func (j *job) SetLogger(logger *logrus.Logger) Job {
	j.log = logger
	return j
}

func (j *job) AddStep(step Step) Job {
	step.SetJobName(j.config.Name)
	step.SetLogger(j.log)

	j.steps = append(j.steps, step)

	return j
}

func (j *job) Run() error {
	if j.config.JobListener != nil {
		j.config.JobListener.Before(j.ctx, j.log)
	}

	for _, step := range j.steps {
		if err := step.WithContext(NewStepContext(j.ctx)).Run(); err != nil {
			return err
		}
	}

	if j.config.JobListener != nil {
		j.config.JobListener.After(j.ctx, j.log)
	}

	return nil
}
