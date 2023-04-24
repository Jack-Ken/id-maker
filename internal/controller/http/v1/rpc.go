package v1

import "github.com/gin-gonic/gin"

type response struct {
	Msg string `json:"msg" example:"message"`
}

type idresponse struct {
	Id int64 `json:"id" example:"message"`
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}

func successResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, response{msg})
}

func successIdResponse(c *gin.Context, code int, id int64) {
	c.JSON(code, idresponse{id})
}
