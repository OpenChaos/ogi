package ogitransformer

import (
	"fmt"

	"github.com/abhishekkr/gol/golenv"

	"github.com/gojekfarm/ogi/instrumentation"

	"github.com/gojekfarm/ogi/logger"
	ogiproducer "github.com/gojekfarm/ogi/producer"
)

type LogTransformer interface {
	Transform(string, ogiproducer.Producer) error
}

type NewLogTransformer func() LogTransformer

var (
	KafkaTopicLabel = golenv.OverrideIfEnv("PRODUCER_KAFKA_TOPIC_LABEL", "app")
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "kubernetes-kafka-log")

	transformerMap = map[string]NewLogTransformer{
		"message-log":          NewMessageLog,
		"kubernetes-kafka-log": NewKubernetesKafkaLog,
	}
)

func validateConfig() {
	var missingVariables string
	if KafkaTopicLabel == "" {
		missingVariables = fmt.Sprintf("%s PRODUCER_KAFKA_TOPIC_LABEL", missingVariables)
	}
	if TransformerType == "" {
		missingVariables = fmt.Sprintf("%s TRANSFORMER_TYPE", missingVariables)
	}

	if missingVariables != "" {
		logger.Fatalf("Missing Env Config:%s", missingVariables)
	}
}

func Transform(producer ogiproducer.Producer, msg string) {
	txn := instrumentation.StartTransaction("transform_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	kafkaLog := transformerMap[TransformerType]()
	if err := kafkaLog.Transform(msg, producer); err != nil {
		logger.Warn(err)
	}
}
