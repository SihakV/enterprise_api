package routers

import (
	"midterm/controllers"
	"midterm/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoute(route *gin.Engine) {

	apiV1 := route.Group("/api/v1")
	{

		//admin routes
		admin := apiV1.Group("/admin")
		admin.POST("login", func(c *gin.Context) {
			controllers.Login(c, 1)
		})

		//employee routes
		employee := apiV1.Group("/employee")
		employee.POST("login", func(c *gin.Context) {
			controllers.Login(c, 2)
		})

		//language routes
		language := apiV1.Group("/language")
		language.GET("/getAll", controllers.GetLanguages)
		language.Use(middleware.AuthMiddleware(1))
		language.POST("/create", controllers.CreateLanguage)
		language.PUT("/edit/:id", controllers.EditLanguage)
		language.DELETE("/delete/:id", controllers.DeleteLanguage)

		//position routes
		position := apiV1.Group("/position")
		position.Use(middleware.AuthMiddleware(1))
		position.POST("/create", controllers.CreatePosition)
		position.GET("/getAll", controllers.GetPositions)
		position.PUT("/edit/:id", controllers.EditPosition)
		position.DELETE("/delete/:id", controllers.DeletePosition)

		//skill routes
		skill := apiV1.Group("/skill")
		skill.GET("/getAll", controllers.GetSkills)
		skill.Use(middleware.AuthMiddleware(1))
		skill.POST("/create", controllers.CreateSkill)
		skill.PUT("/edit/:id", controllers.EditSkill)
		skill.DELETE("/delete/:id", controllers.DeleteSkill)

		//category routes
		category := apiV1.Group("/category")
		category.Use(middleware.AuthMiddleware(1))
		category.GET("/getAll", controllers.GetCategorys)
		category.POST("/create", controllers.CreateCategory)
		category.PUT("/edit/:id", controllers.EditCategory)
		category.DELETE("/delete/:id", controllers.DeleteCategory)

		//user routes
		user := apiV1.Group("/user")
		user.Use(middleware.AuthMiddleware(2))
		user.GET("/getProfile", middleware.AuthMiddleware(2), controllers.GetProfile)
		user.PUT("/editProfile", controllers.EditProfile)
		user.Use(middleware.AuthMiddleware(1))
		user.POST("/create", controllers.CreateUser)
		user.GET("/getAll", controllers.GetUsers)
		user.PUT("/edit/:id", controllers.EditUser)
		user.DELETE("/delete/:id", controllers.DeleteUser)

		//job routes
		job := apiV1.Group("/job")
		job.Use(middleware.AuthMiddleware(1))
		job.POST("/create", controllers.CreateJob)
		job.GET("/getAll", controllers.GetJobs)
		job.DELETE("/delete/:id", controllers.DeleteJob)
		job.PUT("/edit/:id", controllers.EditJob)

		//applicants routes
		applicant := apiV1.Group("/applicant")
		apiV1.POST("/create/:id", controllers.CreateApplicant)
		applicant.Use(middleware.AuthMiddleware(1))
		applicant.GET("/getAll", controllers.GetApplicants)
		applicant.DELETE("/delete/:id", controllers.DeleteApplicant)
		applicant.PUT("/edit/:id", controllers.EditApplicant)
		applicant.GET("/count", controllers.CountApplicants)

		//expense routes
		expense := apiV1.Group("/expense")
		expense.POST("/create", middleware.AuthMiddleware(2), controllers.CreateExpense)
		expense.GET("/getByUser", middleware.AuthMiddleware(2), controllers.GetExpenseEmployee)
		expense.Use(middleware.AuthMiddleware(1))
		expense.GET("/getAll", controllers.GetExpenses)
		expense.PUT("/edit/:id", controllers.EditExpense)
		expense.DELETE("/delete/:id", middleware.AuthMiddleware(2), controllers.DeleteExpense)

		//leave routes
		leave := apiV1.Group("/leave")
		leave.POST("/create", middleware.AuthMiddleware(2), controllers.CreateLeave)
		leave.GET("/getByUser", middleware.AuthMiddleware(2), controllers.GetLeaveEmployee)
		leave.Use(middleware.AuthMiddleware(1))
		leave.GET("/getAll", controllers.GetLeaves)
		leave.PUT("/edit/:id", controllers.EditLeave)
		leave.DELETE("/delete/:id", controllers.DeleteLeave)
	}
}
