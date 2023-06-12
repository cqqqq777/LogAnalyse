namespace go job

include "../common/common.thrift"

struct job_base_info {
    1: i64 job_id
    2: i64 create_time
    3: string job_name
}

struct job_info {
    1: i64 job_id
    2: i64 user_id
    3: i64 create_time
    4: i8 status
    5: string job_name
    6: string file_name
    7: string consequent_file
}

struct create_job_req {
    1: i64 user_id
    2: string file_name
    3: string job_name
    4: string field
}

struct create_job_resp {
    1: common.base_resp base_resp
    2: job_info job_info
}

struct list_job_req {
    1: i64 user_id
}

struct list_job_resp {
    1: common.base_resp base_resp
    2: list<job_base_info> job_list
}

struct get_job_info_req {
    1: i64 user_id
    2: i64 job_id
}

struct get_job_info_resp{
    1: common.base_resp base_resp
    2: job_info job_info
}

struct update_job_status_req{
    1: i64 job_id
    2: i8 status
    3: string consequent_file
}

struct update_job_status_resp {
    1: common.base_resp base_resp
}

service JobService {
    create_job_resp CreateJob(1:create_job_req req)
    list_job_resp ListJob(1:list_job_req req)
    get_job_info_resp GetJobInfo(1:get_job_info_req req)
    update_job_status_resp UpdateJob(1:update_job_status_req req)
}











