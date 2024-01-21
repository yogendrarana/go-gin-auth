package handlers

import "github.com/gin-gonic/gin"

func GetDashboardData(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Dashboard data retrieved successfully.",
		"data": gin.H{
			"dashboard_data": "dashboard data",
		},
	})
}
