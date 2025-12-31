package fileserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

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
	data, _ := l.ctx.Value("fileData").([]byte)
	if len(data) == 0 {
		return nil, errors.New("empty file data")
	}
	if req.Size == 0 {
		req.Size = int64(len(data))
	}
	uid, _ := l.ctx.Value("uid").(uint64)
	rpcResp, err := l.svcCtx.FilesRpc.Upload(l.ctx, &pb.UploadReq{
		Bucket:      req.Bucket,
		ObjectName:  req.Key,
		Data:        data,
		ContentType: req.ContentType,
		Size:        req.Size,
		Hash:        req.Hash,
		Uid:         uid,
	})
	if err != nil {
		return nil, fmt.Errorf("rpc upload error: %w", err)
	}

	return &types.UploadFileResp{
		FileId: rpcResp.FileId,
		Bucket: rpcResp.Bucket,
		Key:    rpcResp.ObjectName,
		Size:   rpcResp.Size,
		ETag:   rpcResp.Etag,
	}, nil
}
