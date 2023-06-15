// Code generated by hertz generator.

package user

import (
	"LogAnalyse/app/service/api/config"
	"LogAnalyse/app/shared/middleware"
	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _apiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.CasbinAuth(config.GlobalCasbinEnforcer),
	}
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.CasbinAuth(config.GlobalCasbinEnforcer),
	}
}
