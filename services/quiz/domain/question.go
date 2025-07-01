package domain

import (
	"github.com/quiz_be/services/core/context"
	model "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper"
)

func (d *domain) GetListQuestion(ctx context.Context, page int, limit int) ([]*model.Question, int32, error) {
	d.logger.DebugCtx(ctx, "GetListQuestion")
	offset := helper.GetOffset(page, limit)
	listQuestion, err := d.questionRepo.Query(ctx).OrderByUpdatedAt(true).
		Offset(offset).
		Limit(limit).
		ResultList()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query failed")
		return nil, 0, err
	}
	count, err := d.questionRepo.Query(ctx).Count()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "count query failed")
	}
	return listQuestion, int32(count), err
}
