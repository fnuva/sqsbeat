package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/fnuva/sqsbeat/config"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// sqsbeat configuration.
type sqsbeat struct {
	done      chan struct{}
	config    config.Config
	client    beat.Client
	sqsClient *sqs.SQS
	queueUrl  *string
}

// New creates an instance of sqsbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	var aws_config aws.Config;
	if (c.Credentials != config.Credentials{} && c.Credentials.SecretKey != "" && c.Credentials.AccessKey != "") {
		aws_config = aws.Config{
			Region:      aws.String(c.Region),
			Credentials: credentials.NewStaticCredentials(c.Credentials.AccessKey, c.Credentials.SecretKey, c.Credentials.Token),
		};
	} else
	{
		aws_config = aws.Config{
			Region: aws.String(c.Region)};
	}
	sess, _ := session.NewSession(&aws_config)
	sqsClient := sqs.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))
	queueInfo, _ := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(c.QueueName),
	})
	bt := &sqsbeat{
		done:      make(chan struct{}),
		config:    c,
		sqsClient: sqsClient,
		queueUrl:  queueInfo.QueueUrl,
	}
	return bt, nil
}

// Run starts sqsbeat.
func (bt *sqsbeat) Run(b *beat.Beat) error {
	logp.Info("sqsbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		bytes, _ := json.Marshal(bt.config);
		logp.Info("config:"+ string(bytes))

		receiveMessage, _ := bt.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            bt.queueUrl,
			MaxNumberOfMessages: aws.Int64(10),
			VisibilityTimeout:   aws.Int64(bt.config.VisibilityTimeout), // 20 seconds
			WaitTimeSeconds:     aws.Int64(bt.config.WaitTimeSeconds),
		})
		for _, entity := range receiveMessage.Messages {
			var event_temp common.MapStr
			json.Unmarshal([]byte(*entity.Body), &event_temp)
			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type": event_temp,
				},
			}
			bt.client.Publish(event)
			logp.Info("Event sent")
			bt.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      bt.queueUrl,
				ReceiptHandle: entity.ReceiptHandle,
			})
			logp.Info("Message Deleted from SQS")

		}
	}
}

// Stop stops sqsbeat.
func (bt *sqsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
