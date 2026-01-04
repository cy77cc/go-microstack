package logic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbortMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAbortMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbortMultipartUploadLogic {
	return &AbortMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AbortMultipartUploadLogic) AbortMultipartUpload(in *pb.AbortMultipartUploadReq) (*pb.AbortMultipartUploadResp, error) {
	if in.UploadId == "" {
		return nil, errors.New("uploadId empty")
	}
	rec, err := l.svcCtx.UploadModel.FindOneByUploadId(l.ctx, in.UploadId)
	if err != nil {
		return nil, err
	}
	// 权限检查
	if in.Uid > 0 && rec.Uploader != in.Uid {
		return nil, xcode.NewErrCodeMsg(xcode.PermissionDenied, "permission denied")
	}

	stor, err := l.svcCtx.Storage.Select(l.ctx, rec.Bucket)
	if err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
	}
	if err = stor.AbortMultipart(l.ctx, rec.Bucket, rec.ObjectName, rec.UploadId); err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.ServerError, "server error")
	}
	// 更新状态
	if err = l.svcCtx.UploadModel.UpdateStatusByUploadId(l.ctx, rec.UploadId, 2, 0); err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
	}
	// 清理分片明细
	_ = l.svcCtx.UploadPartMod.DeleteByUploadId(l.ctx, rec.UploadId)

	return &pb.AbortMultipartUploadResp{UploadId: in.UploadId}, nil
}
