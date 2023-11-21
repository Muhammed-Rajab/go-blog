package routers

import (
	"github.com/Muhammed-Rajab/go-blog/pkg/controllers"
	"github.com/Muhammed-Rajab/go-blog/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func BlogRouter(root *gin.RouterGroup) *gin.RouterGroup {

	router := root.Group("/blog")
	controller := controllers.NewBlogController()
	middlewares := middlewares.NewBlogMiddlewares()
	AuthMiddleware := middlewares.CheckForDashboardKey()
	{
		dashboard := router.Group("/dashboard")
		{
			// Dashboard Home
			dashboard.GET("", AuthMiddleware, controller.DashboardHandler)

			// AUTH
			dashboard.GET("/auth", controller.AuthDashboardHandler)
			dashboard.POST("/auth", controller.AuthDashboard)
			dashboard.GET("/auth/logout", controller.LogoutDashboard)

			// BLOG CREATE
			dashboard.GET("/add", AuthMiddleware, controller.AddBlogHandler)
			dashboard.POST("/add", AuthMiddleware, controller.AddBlog)

			// BLOG UPDATE
			dashboard.GET("/edit/:id", AuthMiddleware, controller.EditBlogHandler)
			dashboard.POST("/edit/:id", AuthMiddleware, controller.EditBlog)

			// IMAGE READ
			dashboard.GET("/images", AuthMiddleware, controller.ImagesHandler)

			// IMAGE CREATE
			dashboard.POST("/images", AuthMiddleware, controller.UploadImages)

			// IMAGE DELETE
			dashboard.POST("/images/:id", AuthMiddleware, controller.DeleteImage)

			// BLOG DELETE
			dashboard.DELETE("/:id", AuthMiddleware, controller.DeleteBlog)

			// BLOG UPDATE
			dashboard.PUT("/:id/toggle_publish", AuthMiddleware, controller.TogglePublishBlog)
			dashboard.PUT("/:id", AuthMiddleware, controller.EditBlog)
		}

		// BLOG HOME
		router.GET("", controller.HomeHandler)

		// BLOG INDIVIDUAL
		router.GET("/:slug", controller.BlogHandler)

	}

	return router
}
