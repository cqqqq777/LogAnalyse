namespace go analyse

include "../common/common.thrift"

struct analyse_req {
    1: i64 user_id
    2:string url
    3:string field
}

service AnalyseService {
    common.nil_resp Analyse(1: analyse_req req)
}