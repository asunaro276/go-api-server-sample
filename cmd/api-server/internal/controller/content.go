package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-api-server-sample/cmd/api-server/internal/application"
)

type ContentController struct {
	createUseCase *application.CreateContentUseCase
	getUseCase    *application.GetContentUseCase
	listUseCase   *application.ListContentsUseCase
	updateUseCase *application.UpdateContentUseCase
	deleteUseCase *application.DeleteContentUseCase
}

func NewContentController(
	createUseCase *application.CreateContentUseCase,
	getUseCase *application.GetContentUseCase,
	listUseCase *application.ListContentsUseCase,
	updateUseCase *application.UpdateContentUseCase,
	deleteUseCase *application.DeleteContentUseCase,
) *ContentController {
	return &ContentController{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

func (ctrl *ContentController) Create(c *gin.Context) {
	var req application.CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なリクエストです",
			"details": err.Error(),
		})
		return
	}

	response, err := ctrl.createUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツの作成に失敗しました",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.Content)
}

func (ctrl *ContentController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なIDです",
			"details": err.Error(),
		})
		return
	}

	req := application.GetContentRequest{ID: uint(id)}
	response, err := ctrl.getUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, application.ErrContentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "指定されたコンテンツが見つかりません",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツの取得に失敗しました",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Content)
}

func (ctrl *ContentController) List(c *gin.Context) {
	var req application.ListContentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なクエリパラメータです",
			"details": err.Error(),
		})
		return
	}

	response, err := ctrl.listUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツ一覧の取得に失敗しました",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *ContentController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なIDです",
			"details": err.Error(),
		})
		return
	}

	var reqBody application.UpdateContentRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なリクエストです",
			"details": err.Error(),
		})
		return
	}

	reqBody.ID = uint(id)
	response, err := ctrl.updateUseCase.Execute(c.Request.Context(), &reqBody)
	if err != nil {
		if errors.Is(err, application.ErrContentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "指定されたコンテンツが見つかりません",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツの更新に失敗しました",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Content)
}

func (ctrl *ContentController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なIDです",
			"details": err.Error(),
		})
		return
	}

	req := application.DeleteContentRequest{ID: uint(id)}
	err = ctrl.deleteUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, application.ErrContentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "指定されたコンテンツが見つかりません",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツの削除に失敗しました",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
