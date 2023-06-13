namespace go analyse

include "../common/common.thrift"

struct analyse_req {
    1: i64 user_id
    2: i64 job_id
    3: string url
    4: string field
    5: string job_name
}

service AnalyseService {
    common.nil_resp Analyse(1: analyse_req req)
}