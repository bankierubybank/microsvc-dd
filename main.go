package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/bankierubybank/microsvc-dd/docs"
	"github.com/bankierubybank/microsvc-dd/routes"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//	@title			Microsvc-dd
//	@version		0.1.0-rc
//	@description	This is a sample API for learning microservices
//	@host			localhost:8080
//	@BasePath		/api/v1
func main() {
	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "OK")
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
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
}

//	@Summary		Get debug information
//	@Description	Get debug information
//	@Tags			debug
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/debug [get]
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
	}
	c.JSON(http.StatusOK, d)
}
