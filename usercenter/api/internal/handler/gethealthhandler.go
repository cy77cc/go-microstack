package handler

import (
	"net/http"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/logic"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetHealthLogic(r.Context(), svcCtx)
		resp, err := l.GetHealth()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
