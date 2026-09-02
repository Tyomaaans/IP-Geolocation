package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ip-geo/internal/ip"
)

func NewUserRouter(
	ipHandler *ip.IpHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")
	{
		ip := v1.Group("/geo/ip")
		{
			ip.GET("", ipHandler.TrackVisit)

			ip.GET("visits/:id", ipHandler.GetIpVisitByID)
			ip.GET("visits/:ip", ipHandler.GetIpVisitByIP)
			ip.GET("/visits", ipHandler.GetIpVisits)
			ip.GET("/visits/today", ipHandler.GetTodayIpVisits)

			ip.GET("/histories/:id", ipHandler.GetIpHistoryByID)
			ip.GET("histories/:ip", ipHandler.GetIpHistoriesByIP)
			ip.GET("/histories", ipHandler.GetIpHistories)
			ip.GET("/histories/today", ipHandler.GetTodayIpHistories)
		}
	}

	return r
}