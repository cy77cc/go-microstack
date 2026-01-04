package files

import (
	"context"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cy77cc/go-microstack/common/xcode"
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

	fileData, ok := l.ctx.Value("fileData").([]byte)
	if !ok {
		return nil, xcode.NewErrCode(xcode.FileUploadFail)
	}

	uid := l.ctx.Value("uid").(uint64)

	// 3. 调用 RPC 服务的 Upload 接口
	res, err := l.svcCtx.FilesRpc.Upload(l.ctx, &pb.UploadReq{
		Bucket: req.Bucket,
		ContentType: req.ContentType,
		Size: int64(len(fileData)),
		Uid: uid,
		Hash: req.Hash,
		Data:   fileData,
	})

	if err != nil {
		return nil, xcode.NewErrCode(xcode.FileUploadFail)
	}


	// 5. 返回响应
	return &types.UploadFileResp{
		FileId: res.FileId,
		Bucket: res.Bucket,
		Key: res.ObjectName,
		Size: res.Size,
		ETag: res.Etag,
	}, nil
}
