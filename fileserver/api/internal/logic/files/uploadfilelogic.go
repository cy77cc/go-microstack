package files

import (
	"context"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(req *types.UploadFileReq) (resp *types.UploadFileResp, err error) {
	// TODO: 实现单文件上传逻辑
	// 1. 读取 HTTP 请求中的文件内容 (MultipartForm)
	// 2. 将文件内容读取为 []byte
	// 3. 调用 RPC 服务的 Upload 接口
	//    res, err := l.svcCtx.FileService.Upload(l.ctx, &fileservice.UploadReq{ ... })
	// 4. 返回 RPC 响应结果

	return
}
