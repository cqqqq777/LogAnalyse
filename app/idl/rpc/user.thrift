namespace go user

include "../common/common.thrift"

struct register_req {
    1: string username
    2: string password
}

struct register_resp{
   1: common.base_resp base_resp
   2: string token
}

struct login_req {
    1: string username
    2: string password
}

struct login_resp{
   1: common.base_resp base_resp
   2: i64 id
   3: string token
}

service UserService{
    register_resp Register(1: register_req req)
    login_resp Login(1: login_req req)
}