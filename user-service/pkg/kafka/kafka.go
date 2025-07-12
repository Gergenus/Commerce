package kafka

import "github.com/IBM/sarama"

func ConnectProducer(brokers []string) sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		panic(err)
	}
	return conn
}
