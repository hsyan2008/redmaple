SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

CREATE DATABASE IF NOT EXISTS `redmaple` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `redmaple`;

DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups` (
  `id` int(10) UNSIGNED NOT NULL,
  `name` varchar(32) NOT NULL,
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户组';

INSERT INTO `groups` (`id`, `name`, `is_deleted`, `updated_at`, `created_at`) VALUES
(1, '管理员', 'N', '2017-04-28 02:41:02', '2017-04-28 02:41:02'),
(2, '开发者', 'N', '2017-04-28 02:41:02', '2017-04-28 02:41:02'),
(3, '测试员', 'N', '2017-04-28 02:41:12', '2017-04-28 02:41:12');

DROP TABLE IF EXISTS `machines`;
CREATE TABLE `machines` (
  `id` int(10) UNSIGNED NOT NULL,
  `name` varchar(32) NOT NULL DEFAULT '',
  `env` tinyint(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '哪个环境，1开发2测试3线上',
  `ip` varchar(16) NOT NULL DEFAULT '',
  `port` varchar(5) NOT NULL DEFAULT '',
  `user` varchar(32) NOT NULL DEFAULT '',
  `auth` text NOT NULL,
  `inner_ip` varchar(16) NOT NULL DEFAULT '',
  `inner_port` varchar(5) NOT NULL DEFAULT '',
  `inner_user` varchar(32) NOT NULL DEFAULT '',
  `inner_auth` text NOT NULL,
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='服务器列表';

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` int(10) UNSIGNED NOT NULL,
  `task_id` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT '0',
  `msg` varchar(255) NOT NULL DEFAULT '',
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='每个任务的操作历史';

DROP TABLE IF EXISTS `projectes`;
CREATE TABLE `projectes` (
  `id` int(11) NOT NULL,
  `name` varchar(32) NOT NULL DEFAULT '',
  `git` varchar(255) NOT NULL DEFAULT '' COMMENT 'git地址',
  `wwwroot` varchar(255) NOT NULL DEFAULT '',
  `dev_wwwroot` varchar(255) NOT NULL DEFAULT '',
  `dev_machine_ids` varchar(64) NOT NULL DEFAULT '',
  `dev_after_release` text NOT NULL,
  `test_wwwroot` varchar(255) NOT NULL DEFAULT '',
  `test_machine_ids` varchar(64) NOT NULL DEFAULT '',
  `test_after_release` text NOT NULL,
  `prod_wwwroot` varchar(255) NOT NULL DEFAULT '',
  `prod_machine_ids` varchar(64) NOT NULL DEFAULT '',
  `prod_after_release` text NOT NULL,
  `is_lock` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '如果有项目已经在开发，则为Y锁定，不让修改git和wwwroot字段',
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='项目表';

DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings` (
  `id` int(10) UNSIGNED NOT NULL,
  `smtp_addr` varchar(128) NOT NULL DEFAULT '',
  `smtp_port` smallint(5) UNSIGNED NOT NULL DEFAULT '0',
  `smtp_ssl` enum('Y','N') NOT NULL DEFAULT 'N',
  `smtp_name` varchar(32) NOT NULL DEFAULT '',
  `smtp_user` varchar(128) NOT NULL DEFAULT '',
  `smtp_pass` varchar(64) NOT NULL DEFAULT '',
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='系统设置';

DROP TABLE IF EXISTS `tasks`;
CREATE TABLE `tasks` (
  `id` int(10) NOT NULL,
  `code` char(20) NOT NULL DEFAULT '' COMMENT '用户id和日期生成的20位唯一标识',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '任务名称',
  `user_id` int(10) NOT NULL DEFAULT '0',
  `review_user_id` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '正在进行review的负责人',
  `test_user_id` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '正在进行测试的负责人',
  `comment` text NOT NULL COMMENT '开发内容',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '状态',
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='开发任务表，开发才能创建';

DROP TABLE IF EXISTS `task_projectes`;
CREATE TABLE `task_projectes` (
  `id` int(10) UNSIGNED NOT NULL,
  `task_id` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `project_id` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '同task表的status',
  `start_commit` char(7) NOT NULL DEFAULT '' COMMENT '创建时间的commit id',
  `end_commit` char(7) NOT NULL DEFAULT '' COMMENT '提测时候的commit id',
  `is_patch` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否已经合并到test分支',
  `is_merge` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否合并到预发布分支',
  `is_finish` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否合并到master分支',
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务和项目关联';

DROP TABLE IF EXISTS `task_reviews`;
CREATE TABLE `task_reviews` (
  `id` int(10) UNSIGNED NOT NULL,
  `task_id` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'review人员',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '同task表的status',
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务和review人员关联';

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) UNSIGNED NOT NULL,
  `name` varchar(32) NOT NULL DEFAULT '',
  `realname` varchar(32) NOT NULL DEFAULT '',
  `password` char(32) NOT NULL DEFAULT '',
  `salt` char(32) NOT NULL DEFAULT '',
  `sign` char(32) NOT NULL DEFAULT '',
  `email` varchar(255) NOT NULL DEFAULT '',
  `group_id` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `is_deleted` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

INSERT INTO `users` (`id`, `name`, `realname`, `password`, `salt`, `sign`, `email`, `group_id`, `is_deleted`, `updated_at`, `created_at`) VALUES
(1, 'admin', '超级管理员', '2c0014d4e6f89e9e0215936b6bb92a45', '830128aaab6ce4305b4805c8bde9759d', '', 'hshengyan@ceibs.edu', 1, 'N', '2017-05-04 01:57:39', '2017-04-28 02:11:02');


ALTER TABLE `groups`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uidx_name` (`name`);

ALTER TABLE `machines`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `messages`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `projectes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uidx_name` (`name`);

ALTER TABLE `settings`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `tasks`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_user_id` (`user_id`);

ALTER TABLE `task_projectes`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_task_id` (`task_id`) USING BTREE;

ALTER TABLE `task_reviews`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uidx_task_id_user_id` (`task_id`,`user_id`) USING BTREE,
  ADD KEY `idx_user_id_status` (`user_id`,`status`);

ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uidx_name` (`name`);


ALTER TABLE `groups`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
ALTER TABLE `machines`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
ALTER TABLE `messages`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
ALTER TABLE `projectes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
ALTER TABLE `settings`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
ALTER TABLE `tasks`
  MODIFY `id` int(10) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
ALTER TABLE `task_projectes`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
ALTER TABLE `task_reviews`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
ALTER TABLE `users`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
