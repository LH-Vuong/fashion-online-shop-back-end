package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"online_fashion_shop/api/common/errs"
)

type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte("0"), nil
	}

	stamp := tt.Unix()
	return []byte(fmt.Sprint(stamp)), nil
}

func (t *Timestamp) UnmarshalParam(src string) error {
	n, err := strconv.Atoi(strings.TrimSpace(src))
	if err != nil {
		return err
	}

	ts := time.Unix(int64(n), 0)
	*t = Timestamp(ts)
	return nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	n, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		return err
	}

	ts := time.Unix(int64(n), 0)
	*t = Timestamp(ts)
	return nil
}

func errCodeToHttpStatusCode(code errs.StatusCode) int {
	switch code {
	case errs.OK:
		return http.StatusOK
	case errs.Cancelled:
		return http.StatusRequestTimeout
	case errs.Unknown:
		return http.StatusInternalServerError
	case errs.InvalidArgument:
		return http.StatusBadRequest
	case errs.DeadlineExceeded:
		return http.StatusRequestTimeout
	case errs.NotFound:
		return http.StatusNotFound
	case errs.AlreadyExists:
		return http.StatusConflict
	case errs.PermissionDenied:
		return http.StatusForbidden
	case errs.ResourceExhausted:
		return http.StatusTooManyRequests
	case errs.FailedPrecondition:
		return http.StatusPreconditionFailed
	case errs.Aborted:
		return http.StatusConflict
	case errs.OutOfRange:
		return http.StatusRequestedRangeNotSatisfiable
	case errs.Unimplemented:
		return http.StatusNotImplemented
	case errs.Internal:
		return http.StatusInternalServerError
	case errs.Unavailable:
		return http.StatusServiceUnavailable
	case errs.DataLoss:
		return http.StatusInternalServerError
	case errs.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

type HttpError struct {
	Message string `json:"message"`
}

type HttpSuccess struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func SendError(c *gin.Context, err error) {
	fmt.Println("Response error", err.Error())
	code := errs.Code(err)
	c.AbortWithStatusJSON(errCodeToHttpStatusCode(code), HttpError{
		Message: err.Error(),
	})
}

func SendOK(c *gin.Context, data interface{}) {
	switch data.(type) {
	case string:
		c.JSON(http.StatusOK, HttpSuccess{Message: data.(string)})
	default:
		c.JSON(http.StatusOK, HttpSuccess{Data: data})
	}
}
