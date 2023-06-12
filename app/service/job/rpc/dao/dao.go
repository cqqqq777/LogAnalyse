package dao

import (
	"LogAnalyse/app/service/job/rpc/model"
	"gorm.io/gorm"
)

type Job struct {
	db *gorm.DB
}

func NewJob(db *gorm.DB) *Job {
	m := db.Migrator()
	if !m.HasTable(&model.Job{}) {
		err := m.CreateTable(&model.Job{})
		if err != nil {
			panic(err)
		}
	}
	return &Job{db: db}
}

func (j *Job) CreateJob(job *model.Job) error {
	return j.db.Create(&job).Error
}

func (j *Job) ListJob(id int64) (list []*model.Job, err error) {
	list = make([]*model.Job, 0)
	err = j.db.Raw("select job_id,create_time,job_name from jobs where user_id = ?", id).Scan(&list).Error
	return
}

func (j *Job) GetJobInfo(id int64) (info *model.Job, err error) {
	info = new(model.Job)
	err = j.db.Raw("select job_id,user_id,status,file_name,job_name,consequent_file,create_time from jobs where job_id = ?", id).Scan(info).Error
	return
}

func (j *Job) UpdateJob(id int64, status int8, consequentFile string) error {
	return j.db.Model(&model.Job{}).Where("job_id = ?", id).Updates(model.Job{Status: status, ConsequentFile: consequentFile}).Error
}
