package users

import (
	"net/http"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/logic/users"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Id uint64 `path:"id"`
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := users.NewDeleteUserLogic(r.Context(), svcCtx)
		err := l.DeleteUser(req.Id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
