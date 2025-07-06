package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/quiz_be/services/core/infra/db"
)

var (
	dbTxCtxKey = "dbTxCtx"
)

func GrpcDatabaseTx(db db.SQL, options ...*sql.TxOptions) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		var result interface{}
		var opts *sql.TxOptions
		if len(options) > 0 {
			opts = options[0]
		}

		tx := db.GetDB().WithContext(ctx).Begin(opts)
		if tx.Error != nil {
			fmt.Println("BEGIN Transaction Error:", tx.Error, info.FullMethod, &req)
			return nil, status.Error(codes.Internal, "")
		}
		defer func() {
			if tx != nil {
				fmt.Println("ROLLBACK Transaction Error:", tx.Error, info.FullMethod, &req)
				tx.Rollback()
			}
		}()
		ctx = context.WithValue(ctx, dbTxCtxKey, tx)
		ctx = DbTxToContext(ctx, db.GetDB().WithContext(ctx))
		result, err = handler(ctx, req)
		if err != nil {
			fmt.Println("ROLLBACK execute error:", err)
			tx.Rollback()
		} else {
			fmt.Println("COMMIT execute success:", info.FullMethod, &req)
			tx.Commit()
		}
		tx = nil // reset tx
		return result, err
	}
}

// DbTxFromContext return db transaction by dbTxCtxKey for each request
func DbTxFromContext(ctx context.Context) *gorm.DB {
	raw, _ := ctx.Value(dbTxCtxKey).(*gorm.DB)
	return raw
}

// DbTxToContext set db transaction into context by dbTxCtxKey
func DbTxToContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbTxCtxKey, db)
}
