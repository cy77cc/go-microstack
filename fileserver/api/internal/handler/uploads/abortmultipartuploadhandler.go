package uploads

import (
	"net/http"

	"github.com/cy77cc/go-microstack/common/pkg/response"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/logic/uploads"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AbortMultipartUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AbortMultipartReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := uploads.NewAbortMultipartUploadLogic(r.Context(), svcCtx)
		err := l.AbortMultipartUpload(&req)
		response.Response(r, w, nil, err)
	}
}
