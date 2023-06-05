package errs

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SuccessStatus = "success"
	FailStatus    = "fail"
	ErrorStatus   = "error"
)

func HandleErrorStatus(ctx *gin.Context, err error, method string) {
	log.Printf("failed in the %s method with the following error: %v", method, err.Error())
	ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": ErrorStatus, "message": err.Error()})
}

func HandleFailStatus(ctx *gin.Context, message string, code int) {
	log.Printf("failed with the following message: %v", message)
	ctx.AbortWithStatusJSON(code, gin.H{"status": FailStatus, "message": message})
}
