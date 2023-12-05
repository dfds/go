package orchestrator

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_Run(t *testing.T) {
	//util.InitializeLogger()
	j := &Job{
		Name:     "dummy",
		Status:   &SyncStatus{active: false},
		Schedule: &Schedule{},
		context:  context.Background(),
		handler: func(ctx context.Context) error {
			return nil
		},
		wg: &sync.WaitGroup{},
	}

	j.Run()
	j.wg.Wait()

	j.Status.SetStatus(true)
	j.handler = func(ctx context.Context) error {
		return errors.New("dummy")
	}

	j.Run()
	j.wg.Wait()
}

func TestNewOrchestrator(t *testing.T) {
	orc := NewOrchestrator(context.Background(), &sync.WaitGroup{})
	assert.NotNil(t, orc)
}

func TestOrchestrator_Init(t *testing.T) {
	orc := NewOrchestrator(context.Background(), &sync.WaitGroup{})
	orc.Init(nil)
}

func TestSyncStatus_InProgress(t *testing.T) {
	ss := &SyncStatus{}
	assert.False(t, ss.InProgress())
}

func TestSyncStatus_SetStatus(t *testing.T) {
	ss := &SyncStatus{}
	assert.False(t, ss.InProgress())
	ss.SetStatus(true)
	assert.True(t, ss.InProgress())
}
