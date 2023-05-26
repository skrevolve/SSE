package routes

import (
	"os"

	"github.com/skrevolve/sse/controllers"

	"github.com/gofiber/fiber/v2"
)

var (
	SECRET =os.Getenv("JWT_SECRET_KEY");
)

func Init(app *fiber.App) {

	// controllers
	NoticeController := &controllers.NoticeController{}

	// api v1 group
	v1 := app.Group("/api/v1")
	{
		// UserController
		notice := v1.Group("notice")
		{
			notice.Get("/urgent", NoticeController.UrgentNotice)
		}

	}
}