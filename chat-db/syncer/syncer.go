package syncer

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/WadeCappa/real_time_chat/chat-db/result"
	"github.com/WadeCappa/real_time_chat/chat-db/store"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/consumer"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
	"github.com/jackc/pgx/v5"
)

var DONE_CODE error = fmt.Errorf("stop-listening")

type Syncer struct {
	PostgresUrl            string
	KafkaHostname          string
	ChannelManagerHostname string
}

func getLastOffset(postgresUrl string) (int64, error) {
	result := store.Call(postgresUrl, func(c *pgx.Conn) result.Result[int64] {
		var lastOffset int64 = sarama.OffsetNewest
		err := c.QueryRow(
			context.Background(),
			"select max(message_id) from messages").Scan(&lastOffset)
		if err != nil {
			return result.Failed[int64](err)
		}
		if lastOffset != sarama.OffsetNewest {
			lastOffset += 1
		}
		return result.Success(lastOffset)
	})

	if result.Err != nil {
		return 0, result.Err
	}

	return *result.Result, nil
}

// we can introduce batching here too to further decrease db load
func (s Syncer) RunSyncer() {

	lastOffset, err := getLastOffset(s.PostgresUrl)
	if err != nil {
		log.Printf("assuming offset start: %v", err)
	}

	for {
		err := consumer.WatchForAllWrites(
			[]string{s.KafkaHostname},
			lastOffset,
			func(e events.Event, m consumer.Metadata) error {
				v := updateDataVisitor{metadata: m, syncer: &s}
				err := e.Visit(&v)
				// there's a race here that will kill the server if our last-offset is out of date
				if err != nil {
					return fmt.Errorf("failed to visit data event %v", err)
				}
				log.Printf("successfully wrote message at offset %d", m.Offset)
				return nil
			})

		log.Printf("stopped listening: %v", err)
	}
}
