namespace go file

include "../common/common.thrift"

struct upload_file_req{
    1: string token(api.query="token")
}

struct download_file_req{
    1: string token(api.query="token")
    2: string FileName(api.query="file_name")
}

struct list_file_req{
    1: string token(api.query="token")
}

service FileServie{
    common.nil_resp UploadFile(1:upload_file_req req)(api.post="/api/file")
    common.nil_resp DownloadFile(1:download_file_req req)(api.get="/api/file")
    common.nil_resp ListFile(1:list_file_req req)(api.get="/api/files")
}