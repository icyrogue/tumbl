package tumbl

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type api struct {
	engine   *gin.Engine
	executor *executor
	puller   *puller
}

func NewAPI(executor *executor, puller *puller) *api {
	return &api{
		puller:   puller,
		executor: executor,
		engine:   gin.Default(),
	}
}

func (a *api) Start() error {
	a.engine.LoadHTMLGlob("./templates/*.html")
	a.engine.GET("/", a.mainPage)
	a.engine.GET("/logs", a.logs)
	a.engine.POST("/pull", a.pull)
	a.engine.POST("/exec/:file", a.exec)
	return a.engine.Run()
}

func (a *api) mainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", gin.H{"Link": a.puller.link})
}

func (a *api) pull(c *gin.Context) {
	link := c.PostForm("link")
	if link != "" {
		err := a.puller.SetURL(link)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}
	files, err := a.puller.Pull()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.HTML(http.StatusOK, "files.html", files)
}

func (a *api) exec(c *gin.Context) {
	file := c.Param("file")
	if err := a.executor.Run(file); err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (a *api) logs(c *gin.Context) {
	logs, err := a.executor.GetLogs()
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.HTML(http.StatusOK, "logs.html", logs)
}
