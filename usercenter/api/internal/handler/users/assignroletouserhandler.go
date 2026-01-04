package users

import (
	"net/http"

	"github.com/cy77cc/go-microstack/common/pkg/response"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/logic/users"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AssignRoleToUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Id uint64 `path:"id"`
			types.AssignRoleBody
		}
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := users.NewAssignRoleToUserLogic(r.Context(), svcCtx)
		err := l.AssignRoleToUser(req.Id, &req.AssignRoleBody)
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			httpx.Ok(w)
		}
	}
}
