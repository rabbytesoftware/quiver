package api

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	"github.com/rabbytesoftware/quiver/internal/repositories"
	"github.com/rabbytesoftware/quiver/internal/usecases"
)

func TestNewAPI(t *testing.T) {
	// Set gin to test mode to avoid debug output
	gin.SetMode(gin.TestMode)

	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)

	api := NewAPI(watcher, usecases)

	if api == nil {
		t.Fatal("NewAPI() returned nil")
	}

	// Test that all dependencies are initialized
	if api.watcher == nil {
		t.Error("API.watcher is nil")
	}

	if api.usecases == nil {
		t.Error("API.usecases is nil")
	}

	if api.router == nil {
		t.Error("API.router is nil")
	}
}

func TestAPIStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)

	api := NewAPI(watcher, usecases)

	// Test that API struct has the expected fields
	if api.watcher != watcher {
		t.Error("API.watcher is not the same instance passed to constructor")
	}

	if api.usecases != usecases {
		t.Error("API.usecases is not the same instance passed to constructor")
	}

	// Test that router is not nil
	if api.router == nil {
		t.Error("API.router is nil")
	}
}

func TestAPI_SetupMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)

	api := NewAPI(watcher, usecases)

	// Test that SetupMiddleware doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetupMiddleware() panicked: %v", r)
		}
	}()

	api.SetupMiddleware()
}

func TestAPI_SetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)

	api := NewAPI(watcher, usecases)

	// Test that SetupRoutes doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetupRoutes() panicked: %v", r)
		}
	}()

	api.SetupRoutes()
}

func TestAPIWithNilDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test with nil usecases (but valid watcher)
	watcher := watcher.NewWatcherService()
	api := NewAPI(watcher, nil)

	// Just verify API was created successfully
	if api == nil {
		t.Error("NewAPI() with nil usecases returned nil")
	}

	// Note: We can't test with nil watcher because it causes a panic
	// when the watcher.GetLevel() method is called in NewAPI()
	// This is expected behavior as the watcher is required for proper operation
}

func TestMultipleAPIInstances(t *testing.T) {
	gin.SetMode(gin.TestMode)

	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)

	api1 := NewAPI(watcher, usecases)
	api2 := NewAPI(watcher, usecases)

	// Different API instances
	if api1 == api2 {
		t.Error("NewAPI() should create new instances each time")
	}

	// But they should reference the same dependencies
	if api1.watcher != api2.watcher {
		t.Error("API instances should reference the same watcher when created with the same watcher")
	}

	if api1.usecases != api2.usecases {
		t.Error("API instances should reference the same usecases when created with the same usecases")
	}

	// Routers should be different instances
	if api1.router == api2.router {
		t.Error("API instances should have different router instances")
	}
}

func TestAPIGinModeConfiguration(t *testing.T) {
	// Test that gin mode is configured based on watcher settings
	gin.SetMode(gin.TestMode) // Reset to test mode

	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)

	// This will configure gin mode based on watcher config
	api := NewAPI(watcher, usecases)

	if api == nil {
		t.Fatal("NewAPI() returned nil")
	}

	// We can't easily test the gin mode without accessing internal gin state
	// But we can verify the API was created successfully
	if api.router == nil {
		t.Error("API router was not initialized")
	}
}

func TestAPIMethodsExist(t *testing.T) {
	gin.SetMode(gin.TestMode)

	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)

	api := NewAPI(watcher, usecases)

	// Test that all expected methods exist by calling them in test mode
	// (Function comparisons to nil are always false in Go)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Method call panicked: %v", r)
		}
	}()

	// Test methods exist by calling them
	api.SetupMiddleware()
	api.SetupRoutes()
}

func TestAPI_SetupMiddleware_Coverage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)
	api := NewAPI(watcher, usecases)

	// Test SetupMiddleware doesn't panic and can be called multiple times
	api.SetupMiddleware()
	api.SetupMiddleware() // Should be safe to call multiple times

	// Verify router still exists after middleware setup
	if api.router == nil {
		t.Error("Router should still exist after SetupMiddleware")
	}
}

func TestAPI_SetupRoutes_Coverage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	watcher := watcher.NewWatcherService()
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := usecases.NewUsecases(repos)
	api := NewAPI(watcher, usecases)

	// Test SetupRoutes doesn't panic
	api.SetupRoutes()

	// Verify router still exists after routes setup
	if api.router == nil {
		t.Error("Router should still exist after SetupRoutes")
	}

	// Test that we can create a second API instance and set up its routes
	api2 := NewAPI(watcher, usecases)
	api2.SetupRoutes()

	if api2.router == nil {
		t.Error("Second API router should exist after SetupRoutes")
	}
}
