package logic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/common/xcode"
	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileUrlLogic {
	return &GetFileUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileUrlLogic) GetFileUrl(in *pb.GetFileUrlReq) (*pb.GetFileUrlResp, error) {
	// 1. 获取文件元数据
	fileInfo, err := l.svcCtx.FileModel.FindOneByFileId(l.ctx, in.FileId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xcode.NewErrCode(xcode.NotFound)
		}
		return nil, err
	}

	// 2. 选择存储后端
	store, err := l.svcCtx.Storage.Select(l.ctx, fileInfo.Bucket)
	if err != nil {
		return nil, err
	}

	// 3. 生成 URL
	url, err := store.Presign(l.ctx, fileInfo.Bucket, fileInfo.ObjectName, in.Expires)
	if err != nil {
		return nil, err
	}

	return &pb.GetFileUrlResp{
		Url: url,
	}, nil
}
