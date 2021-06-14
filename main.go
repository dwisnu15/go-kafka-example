package main

import (
	"context"
	"fmt"
	Kafka "github.com/segmentio/kafka-go"
	"strconv"
	"time"
)

func main() {
	//create context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	go produce(ctx)
	consume(ctx)
}

const (
	topic          = "message-log"
	broker1Address = "localhost:29092"
	//broker2Address = "localhost:9094"

)

func produce(ctx context.Context) {
	//init a counter
	i := 0

	//initialize Kafka writer with broker address, and the topic
	kafkaw := Kafka.Writer{
		//set brokers
		Addr: Kafka.TCP(broker1Address),//broker2Address,

		//set topic
		Topic:        topic,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 15,
	}
	for {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := kafkaw.WriteMessages(ctx, Kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			//create message by assigning byte value, which contains the message
			Value: []byte("Data sent from Go apps" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("tidak dapat menulis pesan :" + err.Error())
		}

		//log message
		fmt.Println("writes: ", i)
		i++
		//sleep for a second
		time.Sleep(time.Second)
	}
}

func consume(ctx context.Context) {
	//initialize Kafka reader with the brokers and topic
	//the groupID identifies the consumer and prevents
	//it from receiving duplicate messages
	kafkar := Kafka.NewReader(Kafka.ReaderConfig{
		Brokers: []string{broker1Address},//broker2Address,

		Topic:   topic,
		GroupID: "group1",
	})
	for {
		//the 'ReadMessage' method blocks until we receive the next event
		msg, err := kafkar.ReadMessage(ctx)
		if err != nil {
			panic("tidak dapat membaca pesan : " + err.Error())
		}

		fmt.Println("terima pesan: ", string(msg.Value))
	}
}
