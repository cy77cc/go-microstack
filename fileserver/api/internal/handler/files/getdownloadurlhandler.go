package files

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/logic/fileserver"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDownloadUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			FileId string `path:"fileId"`
		}
		if err := httpx.Parse(r, &path); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		expires := int64(600)
		if qs := r.URL.Query().Get("expires"); qs != "" {
			if v, err := strconv.ParseInt(qs, 10, 64); err == nil && v > 0 {
				expires = v
			}
		}
		ctx := context.WithValue(r.Context(), "fileId", path.FileId)
		ctx = context.WithValue(ctx, "expires", expires)
		l := fileserver.NewGetDownloadUrlLogic(ctx, svcCtx)
		resp, err := l.GetDownloadUrl()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
