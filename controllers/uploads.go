package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/system"
	"github.com/gin-gonic/gin"
)

//UploadPost handles POST /upload route
func UploadPost(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20) // ~32MB
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, "")
		return
	}
	var uris []string
	fmap := c.Request.MultipartForm.File
	for k := range fmap {
		file, fileHeader, err := c.Request.FormFile(k)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		uri, err := saveFile(fileHeader, file)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		uris = append(uris, uri)
	}
	c.JSON(http.StatusOK, uris)
}

//saveFile saves file to disc and returns its relative uri
func saveFile(fh *multipart.FileHeader, f multipart.File) (string, error) {
	fileExt := filepath.Ext(fh.Filename)
	newName := fmt.Sprint(time.Now().Unix()) + fileExt //unique file name ;D
	uri := "/public/uploads/" + newName
	fullName := filepath.Join(system.UploadsPath(), newName)

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
