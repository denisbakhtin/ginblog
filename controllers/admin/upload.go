package admin

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginbasic/system"
	"github.com/gin-gonic/gin"
)

// POST file
func UploadPost(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20) // ~32MB
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, "")
		return
	}
	uris := make([]string, 0)
	fmap := c.Request.MultipartForm.File
	for k, _ := range fmap {
		file, fileHeader, err := c.Request.FormFile(k)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		if uri, err := saveFile(fileHeader, file); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, "")
			return
		} else {
			uris = append(uris, uri)
		}
	}
	c.JSON(http.StatusOK, uris)
}

func saveFile(fh *multipart.FileHeader, f multipart.File) (string, error) {
	fileExt := filepath.Ext(fh.Filename)
	newName := fmt.Sprint(time.Now().Unix()) + fileExt //unique file name ;D
	uri := "/uploads/" + newName
	fullName := filepath.Join(system.GetConfig().Uploads, newName)

	file, err := os.OpenFile(fullName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, f)
	if err != nil {
		return "", err
	}
	return uri, nil
}
