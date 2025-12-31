package roles

import (
	"net/http"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/logic/roles"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GrantPermissionToRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Id uint64 `path:"id"`
			types.GrantPermissionBody
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := roles.NewGrantPermissionToRoleLogic(r.Context(), svcCtx)
		err := l.GrantPermissionToRole(req.Id, &req.GrantPermissionBody)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
