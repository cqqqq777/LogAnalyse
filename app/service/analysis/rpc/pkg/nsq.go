package pkg

import (
	"fmt"

	"LogAnalyse/app/service/analysis/rpc/config"

	"github.com/bytedance/sonic"
	"github.com/nsqio/go-nsq"
)

type Message struct {
	Status         int8   `json:"status"`
	JobID          int64  `json:"job_id"`
	ConsequentFile string `json:"consequent_file"`
}

type Producer struct {
	Producer *nsq.Producer
}

func NewPublisher() (pro *Producer, err error) {
	pro = new(Producer)
	conf := nsq.NewConfig()
	host := fmt.Sprintf("%s:%d", config.GlobalServerConfig.NsqInfo.Host, config.GlobalServerConfig.NsqInfo.Port)
	pro.Producer, err = nsq.NewProducer(host, conf)
	if err != nil {
		pro.Producer.Stop()
		return
	}
	return
}

func (p *Producer) Produce(msg Message) error {
	body, err := sonic.Marshal(msg)
	if err != nil {
		return fmt.Errorf("cannot marshal: %v", err)
	}
	return p.Producer.Publish(config.GlobalServerConfig.NsqInfo.ProducerTopic, body)
}
