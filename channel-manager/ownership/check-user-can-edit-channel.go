package ownership

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func CheckUserCanEditChannel(channelId, userId int64, c *pgx.Conn) error {
	var fromDb int64
	err := c.QueryRow(context.Background(),
		"select channel_id from channel_members where channel_id=$1 and user_id=$2",
		channelId, userId).Scan(&fromDb)
	if err != nil {
		return fmt.Errorf("could not find ownership record for channel %d and user %d", channelId, userId)
	}

	if fromDb != channelId {
		return fmt.Errorf("user %d does not have perms for channel %d", userId, channelId)
	}

	return nil
}
