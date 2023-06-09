// Code generated by hertz generator.

package file

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

func _uploadfileMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.JwtAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
		middleware.CasbinAuth(config.GlobalCasbinEnforcer),
	}
}

func _downloadfileMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.JwtAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
		middleware.CasbinAuth(config.GlobalCasbinEnforcer),
	}
}

func _listfileMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.JwtAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
		middleware.CasbinAuth(config.GlobalCasbinEnforcer),
	}
}
