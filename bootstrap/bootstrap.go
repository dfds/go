package bootstrap

import (
	"context"
	"github.com/gin-gonic/gin"
	bHttp "go.dfds.cloud/bootstrap/http"
	"go.dfds.cloud/bootstrap/log"
	"go.dfds.cloud/orchestrator"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Manager struct {
	Logger       *zap.Logger
	HttpRouter   *gin.Engine
	HttpServer   *http.Server
	Orchestrator *orchestrator.Orchestrator
}

type ManagerBuilder struct {
	context        context.Context
	enableLogging  bool
	loggingOptions struct {
		enableDebug bool
		logLevel    string
	}
	enableMetrics       bool
	enableHttpRouter    bool
	enableOrchestrator  bool
	orchestratorOptions struct {
		namespace string
	}
}

func (m *ManagerBuilder) AddContext(c context.Context) {
	m.context = c
}

func (m *ManagerBuilder) EnableLogging(enableDebug bool, logLevel string) *ManagerBuilder {
	m.enableLogging = true
	m.loggingOptions.enableDebug = enableDebug
	m.loggingOptions.logLevel = logLevel
	return m
}

func (m *ManagerBuilder) EnableMetrics() *ManagerBuilder {
	m.enableMetrics = true
	return m
}

func (m *ManagerBuilder) EnableHttpRouter() *ManagerBuilder {
	m.enableHttpRouter = true
	return m
}

func (m *ManagerBuilder) EnableOrchestrator(namespace string) *ManagerBuilder {
	m.enableOrchestrator = true
	m.orchestratorOptions.namespace = namespace
	return m
}

func (m *ManagerBuilder) Build() *Manager {
	manager := &Manager{}

	if m.context == nil {
		m.context = context.Background()
	}

	if m.enableLogging {
		log.InitializeLogger(m.loggingOptions.enableDebug, m.loggingOptions.logLevel)
		manager.Logger = log.Logger
	}

	if m.enableHttpRouter {
		router, server := bHttp.NewHttpServer()
		manager.HttpRouter = router
		manager.HttpServer = server
	}

	if m.enableMetrics {
		bHttp.NewMetricsServer()
	}

	if m.enableOrchestrator {
		wg := &sync.WaitGroup{}
		orc := orchestrator.NewOrchestrator(m.context, wg, m.orchestratorOptions.namespace)
		manager.Orchestrator = orc
	}

	return manager
}

func Builder() *ManagerBuilder {
	return &ManagerBuilder{}
}
