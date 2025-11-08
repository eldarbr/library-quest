package quest

import "context"

type Repository interface {
	GetQuest(ctx context.Context, id int) (Quest, error)
}
