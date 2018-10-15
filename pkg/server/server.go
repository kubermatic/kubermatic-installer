package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kubermatic/kubermatic-installer/pkg/assets"
	"github.com/kubermatic/kubermatic-installer/pkg/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var (
	manager *installManager

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// NewServer creates a new echo server that serves the
// static wizard assets and also takes care of receiving
// the manifest and running the installation process,
// providing access to the log as it goes along.
func NewServer(logger *logrus.Logger) *echo.Echo {
	manager = newInstallManager(logger)

	// Echo instance
	e := echo.New()
	e.HideBanner = true

	// configure explicit timeouts
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 10 * time.Second
	e.Server.IdleTimeout = 2 * time.Minute

	// send static assets
	assetServer := http.FileServer(assets.Assets)
	e.GET("/*", func(ctx echo.Context) error {
		assetServer.ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	})

	// perform installations
	e.POST("/install", newInstallHandler(logger))
	e.GET("/logs/:id", newLogsHandler(logger))

	// generate values.yaml
	e.POST("/helm-values", newValuesHandler(logger))

	// during development, the wizard runs on a different port
	// than the installer server
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	return e
}

// newInstallHandler handles POST /install requests and takes
// care of receiving and validating the manifest and then
// kicking off the installation process in a separate goroutine.
func newInstallHandler(logger *logrus.Logger) echo.HandlerFunc {
	type response struct {
		ID string `json:"id"`
	}

	return func(ctx echo.Context) error {
		// get and check manifest first
		manifest, err := getManifest(ctx)
		if err != nil {
			return err
		}

		// start the install goroutine
		id, err := manager.Start(manifest)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to start installation: %v", err))
		}

		// tell the client the ID for fetching the logs
		return ctx.JSON(http.StatusCreated, response{ID: id})
	}
}

// newLogsHandler upgrades the incoming request to a websocket
// and then streams the install logs to the client.
func newLogsHandler(logger *logrus.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")

		// find the installation process
		logs, err := manager.Logs(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		// we are ready to rumble, let's open a websocket
		// to stream the logs and other events to the client
		ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		// stream the log
		for msg := range logs {
			ws.WriteMessage(websocket.TextMessage, msg)
		}

		return nil
	}
}

// newValuesHandler generates a values.yaml based on the submittted
// manifest and returns it to the client.
func newValuesHandler(logger *logrus.Logger) echo.HandlerFunc {
	type response struct {
		Values string `json:"values"`
	}

	return func(ctx echo.Context) error {
		// get and check manifest first
		manifest, err := getManifest(ctx)
		if err != nil {
			return err
		}

		// create kubermatic's values.yaml
		values, err := helm.LoadValuesFromFile("values.example.yaml")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load stock values.yaml.")
		}

		err = values.ApplyManifest(&manifest)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create Helm values.yaml.")
		}

		return ctx.JSON(http.StatusOK, response{
			Values: string(values.YAML()),
		})
	}
}

func getManifest(ctx echo.Context) (manifest.Manifest, error) {
	manifest := manifest.Manifest{}

	manifestYAML := ctx.FormValue("manifest")
	if len(manifestYAML) == 0 {
		return manifest, echo.NewHTTPError(http.StatusBadRequest, "No manifest specified.")
	}

	if err := yaml.Unmarshal([]byte(manifestYAML), &manifest); err != nil {
		return manifest, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Manifest is not valid YAML: %v", err))
	}

	err := manifest.Validate()
	if err != nil {
		return manifest, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Manifest is invalid: %v", err))
	}

	return manifest, nil
}
