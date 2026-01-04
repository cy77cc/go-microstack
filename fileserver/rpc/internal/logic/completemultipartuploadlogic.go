package logic

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"
	"github.com/google/uuid"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CompleteMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteMultipartUploadLogic {
	return &CompleteMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CompleteMultipartUploadLogic) CompleteMultipartUpload(in *pb.CompleteMultipartUploadReq) (*pb.CompleteMultipartUploadResp, error) {
	// 秒传逻辑：如果 Upload 记录里有 Hash，且已存在文件，直接返回
	// 注意：CompleteMultipartUpload 通常意味着已经传完了，但为了幂等性，可以检查
	// 另外，如果 Initiated 时带了 Hash，这里可以先查库
	rec, err := l.svcCtx.UploadModel.FindOneByUploadId(l.ctx, in.UploadId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
	}

	// 权限检查
	if in.Uid > 0 && rec.Uploader != in.Uid {
		return nil, xcode.NewErrCodeMsg(xcode.PermissionDenied, "permission denied")
	}

	// 检查秒传
	if rec.Hash != "" {
		exist, err := l.svcCtx.FileModel.FindOneByHash(l.ctx, rec.Hash)
		if err == nil && exist != nil {
			// 清理本次上传的临时状态
			// 这里有个策略问题：如果秒传成功，是否要 Abort 本次 Multipart？
			// 假设是前端并发检测，如果发现已存在，应该调 Abort。
			// 如果到了 Complete 阶段，说明前端已经传完了分片。
			// 这里我们选择直接返回已存在文件，并标记本次上传为“已完成”（虽然指向旧文件）
			// 或者，我们可以执行合并，但数据库指向旧文件（浪费存储）。
			// 最佳实践：前端在 Initiate 前先 Check Hash。
			// 如果走到这里，假设必须合并。
		}
	}

	stor, err := l.svcCtx.Storage.Select(l.ctx, rec.Bucket)
	if err != nil {
		return nil, err
	}
	var parts []svc.CompletedPart
	for _, p := range in.Parts {
		parts = append(parts, svc.CompletedPart{PartNumber: int(p.PartNumber), ETag: p.Etag})
	}
	if _, err = stor.CompleteMultipart(l.ctx, rec.Bucket, rec.ObjectName, rec.UploadId, parts); err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.ServerError, "server error")
	}
	fileId := uuid.NewString()

	// 开启事务
	err = l.svcCtx.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. 插入文件信息
		_, err = l.svcCtx.FileModel.WithSession(session).Insert(ctx, &model.FileInfo{
			FileId:      fileId,
			FileName:    filepath.Base(rec.ObjectName),
			Bucket:      rec.Bucket,
			ObjectName:  rec.ObjectName,
			Size:        rec.Size,
			ContentType: rec.ContentType,
			Uploader:    rec.Uploader,
			UploadTime:  time.Now().Unix(),
			Hash:        rec.Hash,
			Description: "",
			DeletedTime: 0,
			Status:      1,
		})
		if err != nil {
			return xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
		}

		// 2. 更新上传任务状态
		return l.svcCtx.UploadModel.WithSession(session).UpdateStatusByUploadId(ctx, rec.UploadId, 1, time.Now().Unix())
	})

	if err != nil {
		// 如果数据库更新失败，理论上应该回滚 Storage 的操作（但在 S3 中 Complete 后很难回滚，通常只能由清理任务处理孤儿文件）
		// 或者记录日志告警
		l.Logger.Errorf("CompleteMultipartUpload DB Transact failed: %v, uploadId: %s", err, rec.UploadId)
		return nil, err
	}

	return &pb.CompleteMultipartUploadResp{
		FileId:     fileId,
		Bucket:     rec.Bucket,
		ObjectName: rec.ObjectName,
		Size:       rec.Size,
	}, nil
}
