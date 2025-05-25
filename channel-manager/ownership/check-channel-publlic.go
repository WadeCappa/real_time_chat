package ownership

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func CheckChannelPublic(channelId int64, c *pgx.Conn) error {
	var fromDb int64
	err := c.QueryRow(context.Background(),
		"select id from channels where id=$1 and public = True",
		channelId).Scan(&fromDb)
	if err != nil {
		return fmt.Errorf("failed to find channel %d", channelId)
	}

	if fromDb != channelId {
		return fmt.Errorf("channel %d is not public", channelId)
	}

	return nil
}
