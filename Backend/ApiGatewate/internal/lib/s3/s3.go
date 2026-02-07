package s3

import (
	"fmt"
	"log/slog"

	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/config"
	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/entities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type S3Storage struct {
	*s3.S3
	bucket *string
	log    *slog.Logger
}

func NewS3Storage(l *slog.Logger, cfg config.S3) *S3Storage {
	const op = "S3Storage.NewS3Storage"

	log := l.With(
		slog.String("op", op),
	)

	log.Error("trying to connect to s3")
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.ACCESS_KEY, cfg.SECRET_ACCESS_KEY, ""),
		Endpoint:         aws.String(cfg.ENDPOINT),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Error("failed to connect to s3")
		panic(fmt.Errorf("%s: %w", op, err))
	}
	s3Client := s3.New(newSession)
	log.Error("successfully connected to s3")

	return &S3Storage{s3Client, &cfg.BUCKET_NAME, l}
}

func (s *S3Storage) saveToS3(flsCh chan<- entities.FileResp, errCh chan<- error, file entities.File) {
	uid := uuid.New().String()
	_, err := s.PutObject(&s3.PutObjectInput{
		Body:   file.File,
		Bucket: s.bucket,
		Key:    aws.String(uid),
	})

	if err != nil {
		errCh <- err
	}

	resp := entities.FileResp{
		Filename: file.Filename,
		URL:      fmt.Sprintf("%s/%s/%s", s.Endpoint, *s.bucket, uid),
	}

	flsCh <- resp
}

func (s *S3Storage) Save(files []entities.File) ([]entities.FileResp, error) {
	resp := make([]entities.FileResp, 0, len(files))

	flsChan := make(chan entities.FileResp)
	errChan := make(chan error)

	defer close(flsChan)
	defer close(errChan)

	for _, file := range files {
		go s.saveToS3(flsChan, errChan, file)
	}

	for i := 0; i < len(files); i++ {
		select {
		case url := <-flsChan:
			resp = append(resp, url)
		case err := <-errChan:
			return nil, err
		}
	}

	return resp, nil
}
