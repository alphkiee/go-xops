/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.83.89--测试环境
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : 192.168.83.89:3306
 Source Schema         : xops

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 05/06/2021 14:55:07
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `p_type` varchar(40) DEFAULT NULL,
  `v0` varchar(40) DEFAULT NULL,
  `v1` varchar(40) DEFAULT NULL,
  `v2` varchar(40) DEFAULT NULL,
  `v3` varchar(40) DEFAULT NULL,
  `v4` varchar(40) DEFAULT NULL,
  `v5` varchar(40) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_index` (`p_type`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
BEGIN;
INSERT INTO `casbin_rule` VALUES (2, 'p', 'ceshi', '/v1/dept/create', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (1, 'p', 'ceshi', '/v1/dept/list', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (10, 'p', 'ceshi', '/v1/host/cmd/:ids/:cmds', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (3, 'p', 'ceshi', '/v1/host/create', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (7, 'p', 'ceshi', '/v1/host/delete', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (9, 'p', 'ceshi', '/v1/host/exphost/:ids', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (12, 'p', 'ceshi', '/v1/host/fileDownload', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (11, 'p', 'ceshi', '/v1/host/fileUpload', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (5, 'p', 'ceshi', '/v1/host/info/:id', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (4, 'p', 'ceshi', '/v1/host/list', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (8, 'p', 'ceshi', '/v1/host/test', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (6, 'p', 'ceshi', '/v1/host/update/:id', 'PATCH', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (22, 'p', 'ceshi', '/v1/menu/tree', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (13, 'p', 'ceshi', '/v1/prometheus/host/:key/:job', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (15, 'p', 'ceshi', '/v1/user/create', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (14, 'p', 'ceshi', '/v1/user/info', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (20, 'p', 'ceshi', '/v1/user/info/changePwd', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (21, 'p', 'ceshi', '/v1/user/info/delete', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (19, 'p', 'ceshi', '/v1/user/info/update/:userId', 'PATCH', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (17, 'p', 'ceshi', '/v1/user/info/uploadImg', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (16, 'p', 'ceshi', '/v1/user/list', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (18, 'p', 'ceshi', '/v1/user/update/:userId', 'PATCH', NULL, NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for relation_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `relation_role_menu`;
CREATE TABLE `relation_role_menu` (
  `sys_menu_id` bigint(20) unsigned NOT NULL COMMENT '''自增编号''',
  `sys_role_id` bigint(20) unsigned NOT NULL COMMENT '''自增编号''',
  PRIMARY KEY (`sys_menu_id`,`sys_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tb_cmdb_host
-- ----------------------------
DROP TABLE IF EXISTS `tb_cmdb_host`;
CREATE TABLE `tb_cmdb_host` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime(3) DEFAULT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '''软删除''',
  `host_name` varchar(128) DEFAULT NULL COMMENT '''主机名''',
  `ip` varchar(128) DEFAULT NULL COMMENT '''主机地址''',
  `port` varchar(64) DEFAULT NULL COMMENT '''SSH端口''',
  `os_version` varchar(128) DEFAULT NULL COMMENT '''系统版本''',
  `host_type` varchar(64) DEFAULT NULL COMMENT '''主机类型''',
  `auth_type` longtext COMMENT '''认证类型''',
  `user` varchar(64) DEFAULT NULL COMMENT '''认证用户''',
  `password` varchar(64) DEFAULT NULL COMMENT '''认证密码''',
  `private_key` varchar(128) DEFAULT NULL COMMENT '''秘钥''',
  `key_passphrase` varchar(64) DEFAULT NULL COMMENT '''秘钥''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  PRIMARY KEY (`id`),
  KEY `idx_tb_cmdb_host_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of tb_cmdb_host
-- ----------------------------
BEGIN;
INSERT INTO `tb_cmdb_host` VALUES (1, '2021-01-08 20:02:45.204', '2021-01-08 20:02:45.204', NULL, 'master01', '192.168.88.88', '22', '', '腾讯云', '密码验证', 'root', '123456', '/root/.ssh', '/root/sshkey', 'ceshi');
INSERT INTO `tb_cmdb_host` VALUES (2, '2021-01-08 20:08:56.783', '2021-01-08 20:08:56.783', '2021-01-09 10:40:21.142', 'worker01', '192.168.32.88', '22', '', '阿里云', '密码验证', 'root', '123456', '/root/.ssh', '/root/sshkey', 'ceshi');
INSERT INTO `tb_cmdb_host` VALUES (3, '2021-01-08 20:12:25.185', '2021-01-09 10:28:06.305', '2021-01-09 10:34:27.453', 'worker01', '192.168.214.99', '22', 'Centos*.2.10', '国创云', '密码验证', 'root', '123456', '/root/.ssh', '/root/sshkey', 'ceshi');
INSERT INTO `tb_cmdb_host` VALUES (4, '2021-01-08 20:12:28.037', '2021-01-09 10:28:34.127', NULL, 'worker01', '192.168.83.88', '22', 'Centos*.2.10', 'bb云', '密码验证', 'root', 'Ustc@2020', '/root/.ssh', '/root/sshkey', 'ceshi');
INSERT INTO `tb_cmdb_host` VALUES (5, '2021-01-12 14:44:19.565', '2021-01-12 14:44:19.565', NULL, 'master01', '192.168.83.89', '22', '', '腾讯云', '密码验证', 'root', 'Ustc@2020', '/root/.ssh', NULL, '');
INSERT INTO `tb_cmdb_host` VALUES (6, '2021-01-12 14:44:19.738', '2021-01-12 14:44:19.738', NULL, 'worker01', '192.168.82.90', '22', '', '阿里云', '密码验证', 'root', 'Ustc@2020', '/root/.ssh', NULL, '');
INSERT INTO `tb_cmdb_host` VALUES (7, '2021-01-12 14:44:19.896', '2021-01-12 14:44:19.896', NULL, 'worker01', '192.168.214.99', '22', 'Centos*.2.10', '国创云', '密码验证', 'root', '123456', '/root/.ssh', NULL, '');
INSERT INTO `tb_cmdb_host` VALUES (8, '2021-01-12 14:44:20.039', '2021-01-12 14:44:20.039', NULL, 'worker01', '192.168.80.11', '22', 'Centos*.2.10', 'bb云', '密码验证', 'root', '123456', '/root/.ssh', NULL, '');
INSERT INTO `tb_cmdb_host` VALUES (9, '2021-01-12 14:53:52.728', '2021-01-12 14:53:52.728', NULL, 'master01', '192.168.88.88', '22', '', '腾讯云', '密码验证', 'root', '123456', '/root/.ssh', NULL, 'ceshi');
INSERT INTO `tb_cmdb_host` VALUES (10, '2021-01-12 14:53:52.911', '2021-01-12 14:53:52.911', NULL, 'worker01', '192.168.32.88', '22', '', '阿里云', '密码验证', 'root', '123456', '/root/.ssh', NULL, 'ceshi');
INSERT INTO `tb_cmdb_host` VALUES (11, '2021-01-12 14:53:53.078', '2021-01-12 14:53:53.078', NULL, 'worker01', '192.168.214.99', '22', 'Centos*.2.10', '国创云', '密码验证', 'root', '123456', '/root/.ssh', NULL, 'ceshi');
INSERT INTO `tb_cmdb_host` VALUES (12, '2021-01-12 14:53:53.275', '2021-01-12 14:53:53.275', NULL, 'worker01', '192.168.80.11', '22', 'Centos*.2.10', 'bb云', '密码验证', 'root', '123456', '/root/.ssh', NULL, 'ceshi');
COMMIT;

-- ----------------------------
-- Table structure for tb_sys_api
-- ----------------------------
DROP TABLE IF EXISTS `tb_sys_api`;
CREATE TABLE `tb_sys_api` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime(3) DEFAULT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '''软删除''',
  `name` varchar(64) DEFAULT NULL COMMENT '''接口名称''',
  `method` varchar(64) DEFAULT NULL COMMENT '''请求方式''',
  `path` varchar(128) DEFAULT NULL COMMENT '''访问路径''',
  `category` varchar(128) DEFAULT NULL COMMENT '''所属类别''',
  `desc` varchar(128) DEFAULT NULL COMMENT '''说明''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  PRIMARY KEY (`id`),
  KEY `idx_tb_sys_api_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tb_sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `tb_sys_dept`;
CREATE TABLE `tb_sys_dept` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime(3) DEFAULT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '''软删除''',
  `name` varchar(64) DEFAULT NULL COMMENT '''部门名称''',
  `status` tinyint(1) DEFAULT '1' COMMENT '''状态(正常/禁用, 默认正常)''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `sort` int(3) DEFAULT NULL COMMENT '''排序''',
  `parent_id` bigint(20) unsigned DEFAULT '0' COMMENT '''父级部门(编号为0时表示根)''',
  PRIMARY KEY (`id`),
  KEY `idx_tb_sys_dept_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of tb_sys_dept
-- ----------------------------
BEGIN;
INSERT INTO `tb_sys_dept` VALUES (1, '2021-01-03 14:01:07.000', '2021-01-03 14:01:10.000', NULL, '自动化运维', 1, 'admin', NULL, 0);
INSERT INTO `tb_sys_dept` VALUES (2, '2021-01-07 13:50:27.216', '2021-01-07 13:50:27.216', NULL, 'dddd', 1, 'admin', 0, 0);
INSERT INTO `tb_sys_dept` VALUES (3, '2021-01-07 13:52:33.883', '2021-01-07 13:52:33.883', NULL, 'dddd', 1, 'ceshi', 0, 0);
INSERT INTO `tb_sys_dept` VALUES (4, '2021-01-07 13:54:37.481', '2021-01-07 13:54:37.481', NULL, 'skdjd', 1, 'ceshi', 0, 0);
COMMIT;

-- ----------------------------
-- Table structure for tb_sys_dict
-- ----------------------------
DROP TABLE IF EXISTS `tb_sys_dict`;
CREATE TABLE `tb_sys_dict` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime(3) DEFAULT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '''软删除''',
  `key` varchar(64) DEFAULT NULL COMMENT '''字典Key''',
  `value` varchar(64) DEFAULT NULL COMMENT '''字典Value''',
  `desc` varchar(128) DEFAULT NULL COMMENT '''说明''',
  `status` tinyint(1) DEFAULT '1' COMMENT '''状态(正常/禁用, 默认正常)''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  `parent_id` bigint(20) unsigned DEFAULT '0' COMMENT '''父级字典(编号为0时表示根)''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key` (`key`),
  KEY `idx_tb_sys_dict_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tb_sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `tb_sys_menu`;
CREATE TABLE `tb_sys_menu` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime(3) DEFAULT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '''软删除''',
  `name` varchar(64) DEFAULT NULL COMMENT '''菜单名称''',
  `icon` varchar(64) DEFAULT NULL COMMENT '''菜单图标''',
  `path` varchar(64) DEFAULT NULL COMMENT '''菜单访问路径''',
  `sort` int(3) DEFAULT NULL COMMENT '''菜单顺序(同级菜单, 从0开始, 越小显示越靠前)''',
  `status` tinyint(1) DEFAULT '1' COMMENT '''菜单状态(正常/禁用, 默认正常)''',
  `parent_id` bigint(20) unsigned DEFAULT '0' COMMENT '''父菜单编号(编号为0时表示根菜单)''',
  `creator` varchar(64) DEFAULT NULL COMMENT '''创建人''',
  PRIMARY KEY (`id`),
  KEY `idx_tb_sys_menu_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tb_sys_operlog
-- ----------------------------
DROP TABLE IF EXISTS `tb_sys_operlog`;
CREATE TABLE `tb_sys_operlog` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime(3) DEFAULT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '''软删除''',
  `name` varchar(128) DEFAULT NULL COMMENT '''接口名称''',
  `path` varchar(128) DEFAULT NULL COMMENT '''访问路径''',
  `method` varchar(128) DEFAULT NULL COMMENT '''请求方式''',
  `body` blob COMMENT '''请求主体(通过二进制存储节省空间)''',
  `data` blob COMMENT '''响应数据(通过二进制存储节省空间)''',
  `status` bigint(20) DEFAULT NULL COMMENT '''响应状态码''',
  `username` varchar(128) DEFAULT NULL COMMENT '''用户登录名''',
  `ip` varchar(128) DEFAULT NULL COMMENT '''Ip地址''',
  `ip_location` varchar(128) DEFAULT NULL COMMENT '''Ip所在地''',
  `latency` bigint(20) DEFAULT NULL COMMENT '''请求耗时(ms)''',
  `user_agent` varchar(128) DEFAULT NULL COMMENT '''浏览器标识''',
  PRIMARY KEY (`id`),
  KEY `idx_tb_sys_operlog_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tb_sys_role
-- ----------------------------
DROP TABLE IF EXISTS `tb_sys_role`;
CREATE TABLE `tb_sys_role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime(3) DEFAULT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '''软删除''',
  `name` varchar(128) DEFAULT NULL COMMENT '''角色名称''',
  `keyword` varchar(64) DEFAULT NULL COMMENT '''角色关键词''',
  `desc` varchar(255) DEFAULT NULL COMMENT '''角色说明''',
  `status` tinyint(1) DEFAULT '1' COMMENT '''角色状态(正常/禁用, 默认正常)''',
  `creator` varchar(128) DEFAULT NULL COMMENT '''创建人''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_keyword` (`keyword`),
  KEY `idx_tb_sys_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of tb_sys_role
-- ----------------------------
BEGIN;
INSERT INTO `tb_sys_role` VALUES (1, '2021-01-03 14:01:54.000', '2021-01-03 14:01:56.000', NULL, '测试', 'ceshi', NULL, 1, 'admin');
COMMIT;

-- ----------------------------
-- Table structure for tb_sys_user
-- ----------------------------
DROP TABLE IF EXISTS `tb_sys_user`;
CREATE TABLE `tb_sys_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '''自增编号''',
  `created_at` datetime(3) DEFAULT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '''软删除''',
  `username` varchar(128) DEFAULT NULL COMMENT '''用户名''',
  `password` varchar(128) DEFAULT NULL COMMENT '''密码''',
  `mobile` varchar(128) DEFAULT NULL COMMENT '''手机''',
  `avatar` varchar(128) DEFAULT NULL COMMENT '''头像''',
  `name` varchar(128) DEFAULT NULL COMMENT '''姓名''',
  `email` varchar(128) DEFAULT NULL COMMENT '''邮箱地址''',
  `status` tinyint(1) DEFAULT '1' COMMENT '''用户状态(正常/禁用, 默认正常)''',
  `creator` varchar(128) DEFAULT NULL COMMENT '''创建人''',
  `role_id` bigint(20) unsigned DEFAULT NULL COMMENT '''角色Id外键''',
  `dept_id` bigint(20) unsigned DEFAULT NULL COMMENT '''部门Id外键''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_tb_sys_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of tb_sys_user
-- ----------------------------
BEGIN;
INSERT INTO `tb_sys_user` VALUES (1, '2021-01-03 14:02:17.000', '2021-01-03 14:02:20.000', NULL, 'admin', '$2a$10$.uKzqK7jPyiCJ2.2OsuVrO46jAhUlxiTm4BjzWjz/wYct.AUtts4u', '11111111111', NULL, 'ceshi', '10086@qq.com', 1, 'admin', 1, 1);
INSERT INTO `tb_sys_user` VALUES (2, '2021-02-02 09:48:55.943', '2021-02-02 09:48:55.943', NULL, 'gcuser', '$2a$10$zX4MY5a9.eIohP5SDUQPE.AuY.syXqfQP3dNGT/ICIcLOSTqVlvEG', '', '', 'gc', '', 1, 'ceshi', 1, 0);
INSERT INTO `tb_sys_user` VALUES (3, '2021-02-02 10:35:54.960', '2021-02-02 14:50:15.033', NULL, 'test2021', '$2a$10$Lj0hWgRjMIcB1rSZSoAuderxp3uNtOMqBVuEIz3CwWw6NpOK.bEle', '1577654321', 'string', 'gc2021', 'lijie@qq.com', 1, 'ceshi', 1, 0);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
