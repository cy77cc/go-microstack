package handler

import (
	"net/http"

	"github.com/cy77cc/go-microstack/common/pkg/response"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/logic"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
)

func GetHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetHealthLogic(r.Context(), svcCtx)
		resp, err := l.GetHealth()
		response.Response(r, w, resp, err)
	}
}
