package users

import (
	"net/http"

	"github.com/cy77cc/go-microstack/common/pkg/response"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/logic/users"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := users.NewCreateUserLogic(r.Context(), svcCtx)
		resp, err := l.CreateUser(&req)
		response.Response(r, w, resp, err)
	}
}
