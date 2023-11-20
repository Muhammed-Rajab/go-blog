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
			dashboard.GET("", AuthMiddleware, controller.DashboardHandler)

			dashboard.GET("/auth", controller.AuthDashboardHandler)
			dashboard.POST("/auth", controller.AuthDashboard)
			dashboard.GET("/auth/logout", controller.LogoutDashboard)

			dashboard.GET("/add", AuthMiddleware, controller.AddBlogHandler)
			dashboard.POST("/add", AuthMiddleware, controller.AddBlog)

			dashboard.GET("/edit/:id", AuthMiddleware, controller.EditBlogHandler)
			dashboard.POST("/edit/:id", AuthMiddleware, controller.EditBlog)

			dashboard.GET("/images", AuthMiddleware, controller.ImagesHandler)
			dashboard.POST("/images", AuthMiddleware, controller.UploadImages)
			dashboard.POST("/images/:id", AuthMiddleware, controller.DeleteImage)

			dashboard.DELETE("/:id", AuthMiddleware, controller.DeleteBlog)
			dashboard.PUT("/:id/toggle_publish", AuthMiddleware, controller.TogglePublishBlog)
			dashboard.PUT("/:id", AuthMiddleware, controller.EditBlog)
		}

		router.GET("", controller.HomeHandler)
		router.GET("/:slug", controller.BlogHandler)

	}

	return router
}
