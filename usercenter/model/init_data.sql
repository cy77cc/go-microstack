-- 1. 初始化权限 (Permissions)
-- 用户管理
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `resource`, `action`, `description`, `status`, `create_time`, `update_time`) VALUES
(1, '创建用户', 'user:create', 1, 'users', 'create', '允许创建新用户', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, '删除用户', 'user:delete', 1, 'users', 'delete', '允许删除用户', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, '更新用户', 'user:update', 1, 'users', 'update', '允许更新用户信息', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, '查看用户详情', 'user:get', 1, 'users', 'get', '允许查看用户详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(5, '查看用户列表', 'user:list', 1, 'users', 'list', '允许查看用户列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(6, '分配角色', 'user:assign_role', 1, 'users', 'assign_role', '允许给用户分配角色', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 角色管理
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `resource`, `action`, `description`, `status`, `create_time`, `update_time`) VALUES
(11, '创建角色', 'role:create', 1, 'roles', 'create', '允许创建新角色', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(12, '删除角色', 'role:delete', 1, 'roles', 'delete', '允许删除角色', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(13, '更新角色', 'role:update', 1, 'roles', 'update', '允许更新角色信息', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(14, '查看角色详情', 'role:get', 1, 'roles', 'get', '允许查看角色详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(15, '查看角色列表', 'role:list', 1, 'roles', 'list', '允许查看角色列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(16, '分配权限', 'role:grant_permission', 1, 'roles', 'grant_permission', '允许给角色分配权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 权限管理
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `resource`, `action`, `description`, `status`, `create_time`, `update_time`) VALUES
(21, '创建权限', 'permission:create', 1, 'permissions', 'create', '允许创建新权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(22, '删除权限', 'permission:delete', 1, 'permissions', 'delete', '允许删除权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(23, '更新权限', 'permission:update', 1, 'permissions', 'update', '允许更新权限信息', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(24, '查看权限详情', 'permission:get', 1, 'permissions', 'get', '允许查看权限详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(25, '查看权限列表', 'permission:list', 1, 'permissions', 'list', '允许查看权限列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 2. 初始化角色 (Roles)
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `create_time`, `update_time`) VALUES
(1, '超级管理员', 'admin', '拥有所有权限的超级管理员', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, '普通用户', 'user', '普通注册用户', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 3. 初始化用户 (Users)
-- 密码 hash 对应明文: admin
-- 注意：这里使用的是 bcrypt 生成的 hash，如果你后端使用其他加密方式，请替换
INSERT INTO `users` (`id`, `username`, `password_hash`, `email`, `phone`, `avatar`, `status`, `create_time`, `update_time`) VALUES
(1, 'admin', '3ae17d0c9efcbc645a1ae442f56973119e764ea1f0ee4c420dc43ec3278a0782', 'admin@example.com', '13800000000', '', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 4. 关联用户和角色 (User Roles)
-- 给 admin 用户 (id=1) 分配 admin 角色 (id=1)
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES (1, 1);

-- 5. 关联角色和权限 (Role Permissions)
-- 给 admin 角色 (id=1) 分配所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 1, id FROM `permissions`;