package controllers

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/denisbakhtin/ginblog/config"
	"github.com/gin-gonic/gin"
)

// UploadPost handles POST /upload route
func UploadPost(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20) // ~32MB
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, "")
		return
	}

	mpartFile, mpartHeader, err := c.Request.FormFile("upload")
	if err != nil {
		slog.Error(err.Error())
		c.String(400, err.Error())
		return
	}
	defer mpartFile.Close()
	uri, err := saveFile(mpartHeader, mpartFile)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"uploaded": true, "url": uri})
}

// saveFile saves file to disc and returns its relative uri
func saveFile(fh *multipart.FileHeader, f multipart.File) (string, error) {
	fileExt := filepath.Ext(fh.Filename)
	newName := fmt.Sprint(time.Now().Unix()) + fileExt //unique file name ;D
	uri := "/public/uploads/" + newName
	fullName := filepath.Join(config.UploadsPath(), newName)

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
