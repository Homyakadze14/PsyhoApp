package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/entities"
	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/lib/s3"
	"github.com/gin-gonic/gin"
)

const (
	maxFilesSize    = 10 << 20
	photosArrLenght = 5
)

var (
	ErrHudgeFiles       = errors.New("hudge files")
	ErrWrongContentType = errors.New("content type must be multipart/form-data")
	ErrNoFiles          = errors.New("you must send at least one file")
)

type mediaRoutes struct {
	log *slog.Logger
	s3  *s3.S3Storage
}

func NewMediaRoutes(log *slog.Logger, handler *gin.RouterGroup, s3 *s3.S3Storage) {
	r := &mediaRoutes{
		log: log,
		s3:  s3,
	}

	g := handler.Group("/media")
	{
		g.POST("/upload", r.upload)
	}
}

func (r *mediaRoutes) getFiles(c *gin.Context) ([]entities.File, error) {
	err := c.Request.ParseMultipartForm(maxFilesSize)
	if err != nil {
		return nil, ErrHudgeFiles
	}

	multipartFormData, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	arr := multipartFormData.File["files"]

	if len(arr) == 0 {
		return nil, ErrNoFiles
	}

	files := make([]entities.File, 0, photosArrLenght)
	for _, fileHeader := range arr {
		file := entities.File{
			Filename: fileHeader.Filename,
		}

		f, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()

		file.File = f
		files = append(files, file)
	}

	return files, nil
}

// @Summary     Upload media
// @Description Upload media
// @ID          Upload media
// @Tags  	    media
// @Param 		files formData []file false "files"
// @Accept      mpfd
// @Produce     json
// @Success     200 {object} entities.UploadResponse
// @Failure     400
// @Failure     500
// @Failure     503
// @Router      /media/upload [post]
func (r *mediaRoutes) upload(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		slog.Error(ErrWrongContentType.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrWrongContentType.Error()})
		return
	}

	files, err := r.getFiles(c)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, ErrHudgeFiles) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrHudgeFiles.Error()})
			return
		}
		if errors.Is(err, ErrNoFiles) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrNoFiles.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	fls, err := r.s3.Save(files)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, entities.UploadResponse{Files: fls})
}
