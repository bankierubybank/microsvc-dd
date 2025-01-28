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

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
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

var (
	CIRunNumber = "N/A"
)

func initTracer() (*sdktrace.TracerProvider, error) {
	// Set the OTLP endpoint
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	// Create the OTLP exporter
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(), // Use this if the endpoint is not secured with TLS
	)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	// Create the tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

// @title			Microsvc-dd
// @version		0.1.1-rc
// @description	This is a sample API for learning microservices
// @BasePath		/api/v1
func main() {
	// Initialize OpenTelemetry tracer provider
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

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
	router.Use(otelgin.Middleware("my-server"))

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
// @Tags			debug
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/debug [get]
func GetDebug(c *gin.Context) {
	d := new(debugInfo)

	d.RuntimeInfo.Hostname = os.Getenv("HOSTNAME")

	uname, unameErr := (exec.Command("uname", "-a")).Output()
	if unameErr == nil {
		d.RuntimeInfo.UName = strings.TrimRight(string(uname), "\n")
	}

	d.RuntimeInfo.K8Snode = os.Getenv("NODENAME")
	d.RuntimeInfo.K8Snamespace = os.Getenv("NAMESPACE")

	if info, ok := debug.ReadBuildInfo(); ok {
		d.BuildInfo.GoBuildVersion = info.GoVersion
		for _, setting := range info.Settings {
			if setting.Key == "vcs" {
				d.BuildInfo.VCS = setting.Value
			}
			if setting.Key == "vcs.revision" {
				d.BuildInfo.Commit = setting.Value
				d.BuildInfo.CommitURL = "https://github.com/bankierubybank/microsvc-dd/commit/" + setting.Value
			}
		}
		d.BuildInfo.CIRunNumber = CIRunNumber
	}
	c.JSON(http.StatusOK, d)
}
