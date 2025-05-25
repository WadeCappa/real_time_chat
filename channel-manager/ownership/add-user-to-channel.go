package ownership

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func AddUserToChannel(channelId, userId int64, c *pgx.Conn) error {
	_, err := c.Exec(context.Background(),
		"insert into channel_members (channel_id, user_id) values ($1, $2)",
		channelId,
		userId)
	if err != nil {
		return fmt.Errorf("failed to create new ownership record of channel %d and user %d: %v", channelId, userId, err)
	}
	return err
}
