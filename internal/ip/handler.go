package ip

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ip-geo/pkg"
)

type IpHandler struct {
	ipService  IpService
}

func NewIpHandler(
	ipService  IpService,
) *IpHandler {
	return &IpHandler{
		ipService:  ipService,
	}
}

func httpError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, pkg.ErrNotFound):
		pkg.ErrorResponse(c, http.StatusNotFound, err)
	case errors.Is(err, pkg.ErrInvalidInput):
		pkg.ErrorResponse(c, http.StatusBadRequest, err)
	default:
		pkg.ErrorResponse(c, http.StatusInternalServerError, err)
	}
}

func (h *IpHandler) TrackVisit(c *gin.Context) {
	ip := "103.168.186.82"
	ua := c.GetHeader("User-Agent")

	log.Printf("ip: %s", ip)
	log.Printf("ip: %s", ua)

	if err := h.ipService.TrackVisit(c.Request.Context(), ip, ua); err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "track visit successfully", nil)
}

func (h *IpHandler) GetIpVisitByIP(c *gin.Context) {
	ip := c.Param("ip")

	res, err := h.ipService.GetIpVisitByIP(c.Request.Context(), ip)
	if err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "get ip visit successfully", map[string]interface{}{
		"visit": res,
	})
}

func (h *IpHandler) GetIpVisitByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.ipService.GetIpVisitByID(c.Request.Context(), id)
	if err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "get ip visit successfully", map[string]interface{}{
		"visit": res,
	})
}

func (h *IpHandler) GetIpVisits(c *gin.Context) {
	strPage   := c.DefaultQuery("page", "1")
	strLlimit := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(strLlimit)
	if err != nil {
		limit = 10
	}

	res, err := h.ipService.GetIpVisits(c.Request.Context(), page, limit)
	if err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "get ip visits successfully", res)
}

func (h *IpHandler) GetTodayIpVisits(c *gin.Context) {
	strPage   := c.DefaultQuery("page", "1")
	strLlimit := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(strLlimit)
	if err != nil {
		limit = 10
	}

	res, err := h.ipService.GetTodayIpVisits(c.Request.Context(), page, limit)
	if err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "get today ip visits successfully", res)
}

func (h *IpHandler) GetIpHistoryByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.ipService.GetIpHistoryByID(c.Request.Context(), id)
	if err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "get ip history successfully", map[string]interface{}{
		"history": res,
	})
}

func (h *IpHandler) GetIpHistoriesByIP(c *gin.Context) {
	ip        := c.Param("ip")
	strPage   := c.DefaultQuery("page", "1")
	strLlimit := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(strLlimit)
	if err != nil {
		limit = 10
	}

	res, err := h.ipService.GetIpHistoriesByIP(c.Request.Context(), ip, page, limit)
	if err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "get ip histories successfully", res)
}

func (h *IpHandler) GetIpHistories(c *gin.Context) {
	strPage   := c.DefaultQuery("page", "1")
	strLlimit := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(strLlimit)
	if err != nil {
		limit = 10
	}

	res, err := h.ipService.GetIpHistories(c.Request.Context(), page, limit)
	if err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "get ip histories successfully", res)
}

func (h *IpHandler) GetTodayIpHistories(c *gin.Context) {
	strPage   := c.DefaultQuery("page", "1")
	strLlimit := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(strLlimit)
	if err != nil {
		limit = 10
	}

	res, err := h.ipService.GetTodayIpHistories(c.Request.Context(), page, limit)
	if err != nil {
		httpError(c, err)
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "get today ip histories successfully", res)
}