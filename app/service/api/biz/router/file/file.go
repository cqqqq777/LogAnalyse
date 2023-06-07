// Code generated by hertz generator. DO NOT EDIT.

package file

import (
	file "LogAnalyse/app/service/api/biz/handler/file"
	"github.com/cloudwego/hertz/pkg/app/server"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_api := root.Group("/api", _apiMw()...)
		_api.POST("/file", append(_uploadfileMw(), file.UploadFile)...)
		_api.GET("/file", append(_downloadfileMw(), file.DownloadFile)...)
		_api.GET("/files", append(_listfileMw(), file.ListFile)...)
	}
}