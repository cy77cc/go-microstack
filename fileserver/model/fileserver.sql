CREATE TABLE `file_info`
(
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `file_id`      VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '文件唯一ID (UUID)',
    `file_name`    VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '原始文件名',
    `bucket`       VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '存储桶',
    `object_name`  VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '存储路径/对象名',
    `size`         BIGINT          NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
    `content_type` VARCHAR(128)    NOT NULL DEFAULT '' COMMENT '文件类型（MIME）',
    `uploader`     BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '上传者用户ID',
    `upload_time`  BIGINT          NOT NULL DEFAULT 0 COMMENT '上传时间',
    `hash`         VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '文件哈希值',
    `description`  VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '文件描述',
    `deleted_time` BIGINT          NOT NULL DEFAULT 0 COMMENT '删除时间',
    `status`       TINYINT         NOT NULL DEFAULT 1 COMMENT '状态（1=正常，0=已删除）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_file_id` (`file_id`),
    KEY `idx_uploader` (`uploader`),
    KEY `idx_upload_time` (`upload_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='文件元数据表';
CREATE TABLE `bucket_config`
(
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `bucket`       VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '存储桶',
    `storage_type` TINYINT         NOT NULL DEFAULT 1 COMMENT '存储类型(1=MINIO,2=LOCAL)',
    `endpoint`     VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '存储端点',
    `region`       VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '区域',
    `is_public`    TINYINT         NOT NULL DEFAULT 0 COMMENT '是否公有读',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_bucket` (`bucket`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='存储桶配置';
CREATE TABLE `multipart_upload`
(
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `upload_id`    VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '分片上传ID',
    `bucket`       VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '存储桶',
    `object_name`  VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '对象名',
    `size`         BIGINT          NOT NULL DEFAULT 0 COMMENT '总大小（字节）',
    `content_type` VARCHAR(128)    NOT NULL DEFAULT '' COMMENT '类型',
    `uploader`     BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '发起者用户ID',
    `hash`         VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '哈希',
    `status`       TINYINT         NOT NULL DEFAULT 0 COMMENT '状态(0=进行中,1=完成,2=取消)',
    `create_time`  BIGINT          NOT NULL DEFAULT 0 COMMENT '创建时间',
    `complete_time` BIGINT         NOT NULL DEFAULT 0 COMMENT '完成时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_upload_id` (`upload_id`),
    KEY `idx_bucket_object` (`bucket`, `object_name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='分片上传任务';
CREATE TABLE `multipart_part`
(
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `upload_id`   VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '分片上传ID',
    `part_number` INT             NOT NULL DEFAULT 0 COMMENT '分片序号',
    `etag`        VARCHAR(255)    NOT NULL DEFAULT '' COMMENT 'ETag',
    `size`        BIGINT          NOT NULL DEFAULT 0 COMMENT '分片大小',
    `create_time` BIGINT          NOT NULL DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_part` (`upload_id`, `part_number`),
    KEY `idx_upload_id` (`upload_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='分片上传分片';
