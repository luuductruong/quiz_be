package job

import "github.com/quiz_be/services/core/context"

type Service interface {
	Test(ctx context.Context, name string) (*Test, error)
}
