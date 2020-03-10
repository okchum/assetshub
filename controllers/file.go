package controllers

import (
	"Assets-hub/api"
	"Assets-hub/middlewares/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(path, "./") {
			return nil
		}
		*files = append(*files, path)
		return nil
	}
}

func MultipleUploads(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		api.Out(c, http.StatusBadRequest, fmt.Sprintf("Get form err: %s", err.Error()))
	}
	files := form.File["files"]
	root := config.GetRoot("public_root",true)

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		filename = fmt.Sprintf("%s%s", root, filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			api.Out(c, http.StatusBadRequest, fmt.Sprintf("Upload file err: %s", err.Error()))
		}
	}

	api.Out(c, http.StatusOK, fmt.Sprintf("Uploaded successfully %d files.", len(files)))
}

func GetFilePath(c *gin.Context, publicPath string) string {
	root := config.GetRoot("public_root",false)
	projectPath := c.Query(publicPath)

	if len(strings.TrimSpace(projectPath)) == 0 {
		projectPath = root
	} else if !strings.HasPrefix(projectPath, "/") {
		projectPath = fmt.Sprintf("%s/%s", root, projectPath)
	} else {
		projectPath = fmt.Sprintf("%s%s", root, projectPath)
	}
	return projectPath
}

func ListAll(c *gin.Context) {
	basePath := GetFilePath(c, "path")
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		api.Out(c, http.StatusOK, fmt.Sprintf("%s not exist", basePath))
		return
	}

	var files []string
	err := filepath.Walk(basePath, visit(&files))
	if err != nil {
		panic(err)
	}
	api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s Listing with sub-directories.", basePath), files)
}

func List(c *gin.Context) {
	basePath := GetFilePath(c, "path")
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s is not exist", basePath))
		return
	}

	fileInfo, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, file := range fileInfo {
		fmt.Println(file.Name())
		files = append(files, file.Name())
	}

	api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s Listing.", basePath), files)
}

func Test(c *gin.Context) {
	basePath := GetFilePath(c, "path")

	var files []string
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	api.Out(c, http.StatusOK, fmt.Sprintf("%s", basePath), files)
}

func Move(c *gin.Context) {
	source := GetFilePath(c, "source")
	//destination := GetFilePath(c, "destination")
	if _, err := os.Stat(source); os.IsNotExist(err) {
		api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s is not existed", source))
		return
	}

	if result, err := IsEmptyDir(source); err == nil {
		api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s failed to delete, has files in it. Error: %t", source, result))
		return
	}

	if err := os.Remove(source); err != nil {
		api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s failed to delete, please contact administer. Error: %s", source, err))
		return
	}

	api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s deleted successfully", source))
}

func Remove(c *gin.Context) {
	source := GetFilePath(c, "source")
	force := GetFilePath(c, "force")
	if _, err := os.Stat(source); os.IsNotExist(err) {
		api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s is not existed", source))
		return
	}

	if result, _ := IsEmptyDir(source); !result {
		if err := os.RemoveAll(source); force == "true" && err == nil {
			api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s deleted successfully", source))
		}
		api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s failed to delete, has files in it. Error: %t", source, result))
		return
	}

	if err := os.Remove(source); err != nil {
		api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s failed to delete, please contact administer. Error: %s", source, err))
		return
	}

	api.Out(c, http.StatusOK, fmt.Sprintf("Directory %s deleted successfully", source))
}

func IsEmptyDir(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if _, err = f.Readdir(1); err == io.EOF {
		return true, nil
	}
	return false, err
}
