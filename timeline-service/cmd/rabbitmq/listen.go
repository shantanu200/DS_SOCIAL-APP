package rabbitmq

import (
	"encoding/json"
	"log"
	"timeline/cmd/rabbitmq/event"
)

type Payload struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := DeclareRandomQueue(ch)

	if err != nil {
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(q.Name, s, EXCHANGENAME, false, nil)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload

			_ = json.Unmarshal(d.Body, &payload)

			go HandleEvent(payload)
		}
	}()

	<-forever
	return nil
}

func HandleEvent(payload Payload) {
	switch payload.Event {
	case "FANOUT":
		err := event.UpdateFollowersTimeLine(payload.Data)

		if err != nil {
			log.Printf("[FANOUT ERROR] %s", err.Error())
		}

		log.Println("Fanout service updated")
	case "LIKE":
		err := event.UpdateLikeStatus(payload.Data)

		if err != nil {
			log.Printf("[LIKE ERROR] %s", err.Error())
		}

		log.Println("Like/Dislike service updated")
	case "FOLLOWER":
		err := event.UpdateFollowerUserList(payload.Data)

		if err != nil {
			log.Printf("[FOLLOWER ERROR] %s", err.Error())
		}

		log.Println("Follower service updated")
	case "UPDATE_USER":
        err := event.UpdateUserTweets(payload.Data);

		if err != nil {
			log.Printf("[USER_UPDATE ERROR] %s", err.Error())
		}

		log.Println("Update User service updated")
	
	case "LIKE_REPLY" :
        err := event.UpdateUserReplyLike(payload.Data);

		if err != nil {
			log.Printf("[REPLY_LIKE ERROR] %s", err.Error())
		}

		log.Println("Reply Like service updated")

	default: 
	   log.Fatal("Invalid Operation passed")
	}

}
