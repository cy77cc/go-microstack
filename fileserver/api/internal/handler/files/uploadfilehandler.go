package files

import (
	"context"
	"io"
	"net/http"

	"github.com/cy77cc/go-microstack/common/response"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/logic/fileserver"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadFileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		var data []byte
		// Try to read from multipart form file "file"
		uploadFile, _, err := r.FormFile("file")
		if err == nil && uploadFile != nil {
			defer uploadFile.Close()
			data, err = io.ReadAll(uploadFile)
			if err != nil {
				response.Response(r, w, nil, err)
				return
			}
		} else {
			// Fallback to raw body
			data, err = io.ReadAll(r.Body)
			if err != nil {
				response.Response(r, w, nil, err)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "fileData", data)

		l := fileserver.NewUploadFileLogic(ctx, svcCtx)
		resp, err := l.UploadFile(&req)
		response.Response(r, w, resp, nil)
	}
}
