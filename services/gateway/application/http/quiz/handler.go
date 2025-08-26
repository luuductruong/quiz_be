package quiz

import (
	"encoding/json"
	"github.com/quiz_be/services/gateway/application/http/ws"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net/http"

	"github.com/quiz_be/services/core/application/quiz/dto"
	"github.com/quiz_be/services/core/application/quiz/service"
	"github.com/quiz_be/services/core/context"
	handler "github.com/quiz_be/services/core/http"
)

type HttpHandler interface {
	GetQuizDetail(w http.ResponseWriter, r *http.Request, params map[string]string)
	SubmitAnswer(w http.ResponseWriter, r *http.Request, params map[string]string)
}

func NewHttpHandler(
	bh handler.BaseHandler,
	quizConn *grpc.ClientConn,
	hub *ws.WsHub,
) HttpHandler {
	return &httpHandler{
		BaseHandler: bh,
		service:     service.NewQuizServiceClient(quizConn),
		hub:         hub,
	}
}

type httpHandler struct {
	handler.BaseHandler
	service service.QuizServiceClient
	hub     *ws.WsHub
}

func (h *httpHandler) GetQuizDetail(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fn := func(ctx context.Context, req proto.Message, opts ...grpc.CallOption) (proto.Message, error) {
		reqMapping := req.(*dto.GetQuizDetailReq)
		reqMapping.QuizId = params["quiz_id"]
		res, errResp := h.service.GetQuizDetail(ctx, reqMapping, opts...)
		log.Println("GetQuizDetail return: ", res, "\t err: ", errResp)
		//if errResp == nil && res != nil {
		//	go h.getLeaderboardAndUpdateWS(ctx.Clone(), res.QuizId, 1, 10)
		//}
		return res, errResp
	}
	h.HandleReqWithMd(w, r, &dto.GetQuizDetailReq{}, fn)
}

func (h *httpHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fn := func(ctx context.Context, req proto.Message, opts ...grpc.CallOption) (proto.Message, error) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return nil, err
		}
		reqMapping := req.(*dto.SubmitAnswerReq)
		if err := json.Unmarshal(data, reqMapping); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return nil, err
		}
		res, errResp := h.service.SubmitAnswer(ctx, reqMapping, opts...)
		log.Println("Res: ", res, "\t err: ", errResp)
		//if errResp == nil && res != nil {
		//	go h.getLeaderboardAndUpdateWS(ctx.Clone(), res.Quiz.QuizId, 1, 10)
		//}
		return res, errResp
	}
	h.HandleReqWithMd(w, r, &dto.SubmitAnswerReq{}, fn)
}

func (h *httpHandler) getLeaderboardAndUpdateWS(ctx context.Context, quizID string, page, limit int32) {
	log.Println("getLeaderboardAndUpdateWS\t", quizID, page, limit)
	req := &dto.GetLeaderboardReq{
		QuizId: quizID,
		Limit:  limit,
		Page:   page,
	}
	leaderBoard, err := h.service.GetLeaderboard(ctx, req, grpc.WaitForReady(true))
	if err == nil && leaderBoard != nil {
		go h.hub.PushLeaderboardUpdate(quizID, leaderBoard)
	}
	return
}
