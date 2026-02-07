package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetErrMessages(errs error) error {
	if errs == nil {
		return nil
	}

	newErrMes := ""
	var ve validator.ValidationErrors
	if errors.As(errs, &ve) {
		for _, v := range ve {
			if v.Tag() == "required" {
				newErrMes += fmt.Sprintf("Field %s must be provided;", v.Field())
			}
			if v.Tag() == "email" {
				newErrMes += fmt.Sprintf("Field %s must contains email;", v.Field())
			}
			if v.Tag() == "min" {
				newErrMes += fmt.Sprintf("Minimal lenght for field %s is %v;", v.Field(), v.Param())
			}
			if v.Tag() == "max" {
				newErrMes += fmt.Sprintf("Maximum lenght for field %s is %v;", v.Field(), v.Param())
			}
		}
	} else {
		newErrMes = errs.Error()
	}

	return errors.New(newErrMes)
}

func GetProtoErrWithStatusCode(err error) (int, error) {
	if err == nil {
		return 0, nil
	}

	code := 0
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.InvalidArgument:
			code = http.StatusBadRequest
			err = fmt.Errorf("Invalid argument error: %s", st.Message())
		case codes.NotFound:
			code = http.StatusNotFound
			err = fmt.Errorf("Not found error: %s", st.Message())
		case codes.AlreadyExists:
			code = http.StatusBadRequest
			err = fmt.Errorf("Already exists error: %s", st.Message())
		case codes.Unavailable:
			code = http.StatusServiceUnavailable
			err = fmt.Errorf("Service unavailable")
		case codes.Internal:
			code = http.StatusInternalServerError
			err = fmt.Errorf("Internal server error")
		case codes.Unauthenticated:
			code = http.StatusUnauthorized
			err = fmt.Errorf("Unauthorized: %s", st.Message())
		default:
			code = http.StatusInternalServerError
			err = fmt.Errorf("Unexpected error: %s", st.Message())
		}
	}

	return code, err
}
