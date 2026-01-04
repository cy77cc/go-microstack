package uploads

import (
	"context"
	"io"
	"net/http"

	"github.com/cy77cc/go-microstack/common/pkg/response"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/logic/uploads"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadPartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadPartReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// read part data from body
		data, err := io.ReadAll(r.Body)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		ctx := context.WithValue(r.Context(), "partData", data)
		l := uploads.NewUploadPartLogic(ctx, svcCtx)
		resp, err := l.UploadPart(&req)
		response.Response(r, w, resp, err)
	}
}
