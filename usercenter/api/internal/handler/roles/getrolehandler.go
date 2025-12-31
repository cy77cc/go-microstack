package roles

import (
	"net/http"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/logic/roles"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Id uint64 `path:"id"`
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := roles.NewGetRoleLogic(r.Context(), svcCtx)
		resp, err := l.GetRole(req.Id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
