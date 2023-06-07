namespace go file

include "../common/common.thrift"

struct upload_file_req{
    1: i64 id
}

struct upload_file_resp{
   1: common.base_resp base_resp
   2: string url
}

struct download_file_req{
   1: i64 id
   2: string file_name
}

struct download_file_resp{
   1: common.base_resp base_resp
   2: string url
}

struct file_info{
    1: string name
    2: i64 size
    3: i64 last_modify_time
}

struct list_file_req{
    1: i64 id
}

struct list_file_resp{
    1: common.base_resp base_resp
    2: list<file_info> file_list
}

service FileService{
    upload_file_resp UploadFile(1:upload_file_req req)
    download_file_resp DownloadFile(1:download_file_req req)
    list_file_resp ListFile(1:list_file_req req)
}