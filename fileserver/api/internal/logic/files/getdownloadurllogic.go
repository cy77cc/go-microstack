package files

import (
	"context"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDownloadUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDownloadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDownloadUrlLogic {
	return &GetDownloadUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDownloadUrlLogic) GetDownloadUrl() (resp *types.PresignResp, err error) {
	fileId, _ := l.ctx.Value("fileId").(string)
	if fileId == "" {
		return nil, xcode.NewErrCodeMsg(xcode.ErrInvalidParam, "fileId is empty")
	}
	expires, _ := l.ctx.Value("expires").(int64)
	if expires <= 0 {
		expires = 600
	}
	rpcResp, err := l.svcCtx.FilesRpc.GetFileUrl(l.ctx, &pb.GetFileUrlReq{
		FileId:  fileId,
		Expires: expires,
		Uid:     0,
	})
	if err != nil {
		return nil, err
	}

	return &types.PresignResp{
		Url:    rpcResp.Url,
		Expire: expires,
	}, nil
}
