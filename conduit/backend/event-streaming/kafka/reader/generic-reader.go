package reader

import (
	"backend/event-streaming/kafka/messages/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

//TODO figure out if I am going to finish this cause I think I am going to need to need to make a generic REST event Mux which
// seems like a lot of work, or I could function carry and make every service that wants to use this deal with then
// so in short I don't think all that work is worth do either of those ideas for something I don't think I am going to use
// Look to Dev Notes #220718 for information

const em = "-event-messaging"

//var errGenerticKafkaReaderMessage = errors.New("[ERROR] [EVENT] [READER]")

func generticKafkaReader(logger *log.Logger, consumerGroupReader kafka.Reader) error {
	var errReadingEventMessage = fmt.Errorf("[ERROR] [EVENT] [READER] [%v%v]: Unexpected Error Occuied while attemping to Read Incoming Messages. Restart is Needed. | ", consumerGroupReader.Config().Topic, em)

	logger.Println("Starting Kafka Reader for topic: ", consumerGroupReader.Config().Topic)

	for {
		readMessage, err := consumerGroupReader.ReadMessage(context.Background())
		if err != nil {
			logger.Println(errReadingEventMessage)
			break
		}
		var eventStructs model.EventMessage
		err = json.Unmarshal(readMessage.Value, &eventStructs)
		if err != nil {
			//logger.Println(ErrJSONUnmarshalling, err)
		}
		logger.Println("%+v\n", &eventStructs)
		//Processing read message into what should be done.
		//serverEventHandler.EventMux(eventStructs)
	}
	return nil
}
