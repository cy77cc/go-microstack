package files

import (
	"context"
	"net/http"

	"github.com/cy77cc/go-microstack/common/response"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/logic/fileserver"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFileMetaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			FileId string `path:"fileId"`
		}
		if err := httpx.Parse(r, &path); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		ctx := context.WithValue(r.Context(), "fileId", path.FileId)
		l := fileserver.NewGetFileMetaLogic(ctx, svcCtx)
		resp, err := l.GetFileMeta()
		response.Response(r, w, resp, err)
	}
}
