package handler

import (
	"net/http"

	"github.com/summerKK/mall/service/product/api/internal/logic"
	"github.com/summerKK/mall/service/product/api/internal/svc"
	"github.com/summerKK/mall/service/product/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RemoveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewRemoveLogic(r.Context(), svcCtx)
		resp, err := l.Remove(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
