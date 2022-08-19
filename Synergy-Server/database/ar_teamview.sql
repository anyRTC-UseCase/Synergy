DROP DATABASE IF EXISTS `ar_teamview`;
CREATE DATABASE IF NOT EXISTS `ar_teamview` DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
USE `ar_teamview`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for room_info
-- ----------------------------
DROP TABLE IF EXISTS `room_info`;
CREATE TABLE `room_info`  (
  `roomid` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间id',
  `r_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间名',
  `r_hostid` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房主id',
  `r_state` tinyint(4) NOT NULL DEFAULT 1 COMMENT '房间状态(1:结束,2:进行中,3:转码中)',
  `r_ts` bigint(20) NOT NULL DEFAULT 0 COMMENT '房间创建时间戳(单位:秒)',
  `r_vod_uid` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间录制在频道内使用的UID',
  `r_vod_start_ts` bigint(20) NOT NULL DEFAULT 0 COMMENT '房间录制开始时间(单位:毫秒)',
  `r_vod_stop_ts` bigint(20) NOT NULL DEFAULT 0 COMMENT '房间录制结束时间(单位:毫秒)',
  `r_vod_file_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间录像文件url',
  `r_vod_resource_id` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间录像的resource ID',
  `r_vod_sid` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间录像的录制 ID',
  PRIMARY KEY (`roomid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '房间信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`  (
  `uid` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户id',
  `u_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `u_workname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户工种',
  `u_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '用户类型(1:智能终端,2:专家,3:管理员)',
  `u_ts` bigint(20) NOT NULL DEFAULT 0 COMMENT '用户登录时间戳(单位:秒)',
  PRIMARY KEY (`uid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_join_info
-- ----------------------------
DROP TABLE IF EXISTS `user_join_info`;
CREATE TABLE `user_join_info`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户在房间的自增长id',
  `ar_uid` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户id',
  `ar_user_role` tinyint(4) NOT NULL COMMENT '用户角色(1:主播,2:观众)',
  `ar_roomid` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户在房间的roomid',
  `ar_join_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '用户进入房间的时间戳（单位:秒）',
  `ar_leave_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '用户离开房间的时间戳（单位:秒）',
  `ar_ts` bigint(20) NOT NULL DEFAULT 0 COMMENT '入库时间戳(单位:秒)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户进出房间信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_online_info
-- ----------------------------
DROP TABLE IF EXISTS `user_online_info`;
CREATE TABLE `user_online_info`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增长id',
  `ar_uid` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户id',
  `ar_roomid` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户所在房间id',
  `ar_opt_ts` bigint(20) NOT NULL COMMENT '心跳包时间戳(单位:秒)',
  `ar_ts` bigint(20) NOT NULL DEFAULT 0 COMMENT '入库时间戳(单位:秒)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户在线心跳包信息表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;


2.0数据库更新
房间状态添加通话结束状态
ALTER TABLE `ar_teamview`.`room_info` 
MODIFY COLUMN `r_state` tinyint(4) NOT NULL DEFAULT 4 COMMENT '房间状态(1:录像结束,2:进行中,3:转码中,4:通话结束)' AFTER `r_hostid`;
