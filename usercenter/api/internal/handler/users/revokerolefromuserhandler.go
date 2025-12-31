package users

import (
	"net/http"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/logic/users"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RevokeRoleFromUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Id uint64 `path:"id"`
			types.RevokeRoleBody
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := users.NewRevokeRoleFromUserLogic(r.Context(), svcCtx)
		err := l.RevokeRoleFromUser(req.Id, &req.RevokeRoleBody)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
