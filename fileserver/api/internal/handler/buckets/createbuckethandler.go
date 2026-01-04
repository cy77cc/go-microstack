package buckets

import (
	"net/http"

	"github.com/cy77cc/go-microstack/common/pkg/response"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/logic/buckets"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateBucketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BucketReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := buckets.NewCreateBucketLogic(r.Context(), svcCtx)
		resp, err := l.CreateBucket(&req)
		response.Response(r, w, resp, err)
	}
}
