namespace go user

include "../common/common.thrift"

struct register_req{
    1: string username(api.raw="username", api.vd="len($)>0 && len($)<33")
    2: string password(api.raw="password", api.vd="len($)>0 && len($)<33")
}

struct login_req{
    1: string username(api.raw="username", api.vd="len($)>0 && len($)<33")
    2: string password(api.raw="password", api.vd="len($)>0 && len($)<33")
}

service UserService{
    common.nil_resp Register(1:register_req req)(api.post="/api/user/register")
    common.nil_resp Login(1:login_req req)(api.post="/api/user/login")
}