package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/bankierubybank/microsvc-dd/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1 := router.Group("/api/v1")
	{

		debugRouter := v1.Group("/debug")
		{
			debugRouter.GET("", GetDebug)
		}
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

// @BasePath	/api/v1
// @Summary		Get debug information
// @Schemes
// @Description	Get debug information
// @Tags		debug
// @Accept		json
// @Produce		json
// @Success		200
// @Router		/debug [get]
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
