package files

import (
	"context"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileMetaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileMetaLogic {
	return &GetFileMetaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileMetaLogic) GetFileMeta() (resp *types.FileMetaResp, err error) {
	fileId, _ := l.ctx.Value("fileId").(string)
	if fileId == "" {
		return nil, xcode.NewErrCodeMsg(xcode.ErrInvalidParam, "fileId empty")
	}
	rpcResp, err := l.svcCtx.FilesRpc.GetFileMeta(l.ctx, &pb.GetFileMetaReq{
		FileId: fileId,
	})
	if err != nil {
		return nil, err
	}
	m := rpcResp.Meta
	if m == nil {
		return nil, xcode.NewErrCodeMsg(xcode.NotFound, "meta not found")
	}

	return &types.FileMetaResp{
		FileId:      m.FileId,
		FileName:    m.FileName,
		Bucket:      m.Bucket,
		Key:         m.ObjectName,
		Size:        m.Size,
		ContentType: m.ContentType,
		Uploader:    m.Uploader,
		UploadTime:  m.UploadTime,
		Hash:        m.Hash,
		Status:      m.Status,
	}, nil
}
