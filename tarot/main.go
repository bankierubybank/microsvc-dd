package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/bankierubybank/microsvc-dd/tarot/docs"
	"github.com/bankierubybank/microsvc-dd/tarot/routes"
	"google.golang.org/grpc/credentials"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type debugInfo struct {
	RuntimeInfo runtimeInfo `json:"runtimeInfo"`
	BuildInfo   buildInfo   `json:"buildInfo"`
}

type runtimeInfo struct {
	Hostname     string `json:"hostname"`
	UName        string `json:"uname"`
	K8Snode      string `json:"k8snode"`
	K8Snamespace string `json:"k8snamespace"`
}
type buildInfo struct {
	GoBuildVersion string `json:"gobuildversion"`
	VCS            string `json:"vcs"`
	Commit         string `json:"commit"`
	CommitURL      string `json:"commiturl"`
	CIRunNumber    string `json:"cirunnumber"`
}

type errorInfo struct {
	ErrorCode int    `json:"errorcode"`
	ErrorMesg string `json:"errormesg"`
}

var (
	CIRunNumber  = "N/A"
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

// Ref: https://opentelemetry.io/docs/languages/go/getting-started/
func initTracer() func(context.Context) error {

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

// @title		Microsvc-dd
// @version		0.1.1-rc
// @description	This is a sample API for learning microservices
// @BasePath	/api/v1
func main() {
	// Initialize OpenTelemetry tracer provider
	cleanup := initTracer()
	defer cleanup(context.Background())

	router := gin.Default()

	// Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// OpenTelemetry middleware
	router.Use(otelgin.Middleware(serviceName))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{

		debugRouter := v1.Group("/debug")
		{
			debugRouter.GET("", GetDebug)
		}
		routes.Tarots(v1.Group("/tarots"))
	}
	router.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))) // listen and serve on 0.0.0.0:$PORT
}

// @Summary		Get debug information
// @Description	Get debug information
// @Tags		debug
// @Accept		json
// @Produce		json
// @Success		200 {object} debugInfo
// @Failure		500 {object} errorInfo
// @Router		/debug [get]
func GetDebug(c *gin.Context) {
	d, err := createDebugInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorInfo{
			ErrorCode: http.StatusInternalServerError,
			ErrorMesg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, d)
}

// createDebugInfo creates and populates a debugInfo struct.
func createDebugInfo() (*debugInfo, error) {
	d := &debugInfo{
		RuntimeInfo: runtimeInfo{
			Hostname:     os.Getenv("HOSTNAME"),
			K8Snode:      os.Getenv("NODENAME"),
			K8Snamespace: os.Getenv("NAMESPACE"),
		},
	}

	uname, err := exec.Command("uname", "-a").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get uname information: %w", err)
	}
	d.RuntimeInfo.UName = strings.TrimRight(string(uname), "\n")

	if info, ok := debug.ReadBuildInfo(); ok {
		d.BuildInfo.GoBuildVersion = info.GoVersion
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs":
				d.BuildInfo.VCS = setting.Value
			case "commit":
				d.BuildInfo.Commit = setting.Value
			case "commiturl":
				d.BuildInfo.CommitURL = setting.Value
			}
		}
	} else {
		return nil, fmt.Errorf("failed to read build information")
	}

	return d, nil
}
