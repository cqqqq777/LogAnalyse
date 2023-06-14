namespace go job

include "../common/common.thrift"

struct create_job_req {
    1: string token(api.query="token")
    2: string file_name(api.raw="file_name")
    3: string job_name(api.raw="job_name")
    4: string field(api.raw="field")
}

struct list_job_req {
    1: string token(api.query="token")
}

struct get_job_info_req{
    1: string token(api.query="token")
    2: i64 job_id(api.query="job_id")
}

service JobService{
   common.nil_resp CreateJob(1: create_job_req req)(api.post="/api/job")
   common.nil_resp ListJob(1: list_job_req req)(api.get="/api/jobs")
   common.nil_resp GetJobInfo(1: get_job_info_req req)(api.get="/api/job")
}
