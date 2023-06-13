package pkg

import (
	"LogAnalyse/app/service/api/config"
	"LogAnalyse/app/shared/errz"
	"LogAnalyse/app/shared/kitex_gen/job"
	"LogAnalyse/app/shared/kitex_gen/job/jobservice"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/nsqio/go-nsq"
)

type Consumer struct {
	Consumer *nsq.Consumer
}

type Message struct {
	Status         int8   `json:"status"`
	JobID          int64  `json:"job_id"`
	ConsequentFile string `json:"consequent_file"`
}

func NewConsumer() (consumer *Consumer, err error) {
	consumer = new(Consumer)
	conf := nsq.NewConfig()
	consumer.Consumer, err = nsq.NewConsumer(config.GlobalServerConfig.NsqInfo.ConsumerTopic, config.GlobalServerConfig.NsqInfo.Channel, conf)
	return
}

func handleMsg(client jobservice.Client) func(*nsq.Message) error {
	return func(msg *nsq.Message) error {
		consumerMsg := new(Message)
		err := sonic.Unmarshal(msg.Body, consumerMsg)
		if err != nil {
			return err
		}
		resp, err := client.UpdateJob(context.Background(), &job.UpdateJobStatusReq{JobId: consumerMsg.JobID,
			Status:         consumerMsg.Status,
			ConsequentFile: consumerMsg.ConsequentFile},
		)
		if err != nil {
			return err
		}
		if resp.BaseResp.Code != errz.Success {
			return fmt.Errorf(resp.BaseResp.Msg)
		}
		return nil
	}
}

func (c *Consumer) Consume(client jobservice.Client) error {
	host := fmt.Sprintf("%s:%d", config.GlobalServerConfig.NsqInfo.Host, config.GlobalServerConfig.NsqInfo.Port)
	c.Consumer.AddHandler(nsq.HandlerFunc(handleMsg(client)))
	err := c.Consumer.ConnectToNSQD(host)
	if err != nil {
		return err
	}
	select {}
}
