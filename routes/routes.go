package routes

import (
	adminController "armandwipangestu/gis-api/controllers/admin"
	authController "armandwipangestu/gis-api/controllers/auth"
	"armandwipangestu/gis-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Initialize gin
	router := gin.Default()

	// Auth routes (no auth required)
	auth := router.Group("/api")
	{
		auth.POST("/login", authController.Login)
	}

	// Protected routes (require authentication)
	protected := router.Group("/api/admin")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Dashboard routes
		protected.GET("/dashboard", middlewares.Permission("dashboard-index"), adminController.Dashboard)

		// Permission routes
		protected.GET("/permissions", middlewares.Permission("permissions-index"), adminController.FindPermissions)
		protected.POST("/permissions", middlewares.Permission("permissions-create"), adminController.CreatePermission)
		protected.GET("/permissions/:id", middlewares.Permission("permissions-show"), adminController.FindPermissionById)
		protected.PUT("/permissions/:id", middlewares.Permission("permissions-update"), adminController.UpdatePermission)
		protected.DELETE("/permissions/:id", middlewares.Permission("permissions-delete"), adminController.DeletePermission)
		protected.GET("/permissions/all", middlewares.Permission("permissions-index"), adminController.FindAllPermissions)

		// Role routes
		protected.GET("/roles", middlewares.Permission("roles-index"), adminController.FindRoles)
		protected.POST("/roles", middlewares.Permission("roles-create"), adminController.CreateRole)
		protected.GET("/roles/:id", middlewares.Permission("roles-show"), adminController.FindByRoleId)
		protected.PUT("/roles/:id", middlewares.Permission("roles-update"), adminController.UpdateRole)
		protected.DELETE("/roles/:id", middlewares.Permission("roles-delete"), adminController.DeleteRole)
		protected.GET("/roles/all", middlewares.Permission("roles-index"), adminController.FindAllRoles)

		// User routes
		protected.GET("/users", middlewares.Permission("users-index"), adminController.FindUsers)
	}

	return router
}