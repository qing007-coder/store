package tools

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	MessageBadRequest                 = "Bad Request"
	MessageBadRequestInsufficientData = "Bad Request insufficient data"
	MessageAlreadyVerified            = "Already Verified."
	MessageVerificationAlreadyFailed  = "Verification already Failed."
	MessageVerificationFailed         = "Verification Failed."
	MessagePasskeyVerificationSuccess = "Passkey Authorised."
	MessageInvalidVerificationCode    = "Verification Code Invalid."
	MessageExpiredVerificationCode    = "Verification Code Expired."
	MessageInvalidBody                = "Invalid Body"
	MessageInvalidEmailAddress        = "Invalid Email Address Please use valid Email."
	MessageInvalidPublicKey           = "Invalid PublicKey Please use valid PublicKey."
)

func BadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusBadRequest,
		"message": message,
	})

	ctx.Abort()
}

func InternalServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusInternalServerError,
		"message": "internal server error",
	})

	ctx.Abort()
}

func StatusOK(ctx *gin.Context, data interface{}, message string) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    data,
		"message": message,
	})

	ctx.Abort()
}

func UnauthorizationRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusUnauthorized,
		"message": message,
	})

	ctx.Abort()
}
