package back

import (
	"fmt"
	"net/http"
	"os"

	// "path/filepath"
	"strings"

	"io"

	"github.com/gin-gonic/gin"
)

// var dir = "C:\\Users\\Sergio\\Pictures\\Ahri"

// var dir = "D:\\"
var dir = ""

var abs_dir = dir

var run_dir string

func Run(port int, path string) {
	run_dir = path
	if run_dir == "" {
		run_dir = "."
	}
	fmt.Println("Running server on directory: ", run_dir)

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.Discard

	createPreviews(dir)

	dir = os.Getenv("HOME")
	abs_dir = dir

	// set up cert

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		// set cache control to 1 minute
		// c.Header("Cache-Control", "max-age=60")
	})

	r.POST("/gohome", func(c *gin.Context) {
		// fmt.Println(c.Request.Body)
		dir = ""
		abs_dir = dir
		c.Redirect(http.StatusFound, "/")
	})

	r.POST("/setdir/*root", func(c *gin.Context) {
		dir = c.Param("root")
		dir = strings.ReplaceAll(c.Param("root"), "/", string(os.PathSeparator))
		dir = strings.TrimPrefix(dir, string(os.PathSeparator))
		abs_dir = dir
		grid(c)
	})

	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		fmt.Println("Searching for: ", query)
	})

	r.StaticFS("/static", http.Dir(run_dir+string(os.PathSeparator)+"static"))

	r.GET("/files/*path", func(c *gin.Context) {
		file_path := c.Param("path")
		abs_path := abs_dir + string(os.PathSeparator) + file_path
		files(abs_path, c)
		// c.File(abs_path)
	})

	r.StaticFS("/previews", http.Dir(dir+string(os.PathSeparator)+"previews"))

	r.GET("/", folder)

	r.GET("/:folder/*path", folder)

	r.NoRoute(func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/previews") {
			previews(c)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	r.GET("/grid/:dir/*deer", grid)
	r.GET("/grid/:dir", grid)
	r.GET("/grid", grid)

	// middleware for cache control

	fmt.Println("Starting server on port ", port)
	var err error
	err = r.Run("localhost:" + fmt.Sprint(port))
	if err != nil {
		fmt.Println("Error starting server: ", err)
		os.Exit(1)
	}
}
