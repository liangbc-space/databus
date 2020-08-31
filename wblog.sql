/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 50649
 Source Host           : 127.0.0.1:3306
 Source Schema         : wblog

 Target Server Type    : MySQL
 Target Server Version : 50649
 File Encoding         : 65001

 Date: 18/08/2020 09:51:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for wb_user
-- ----------------------------
DROP TABLE IF EXISTS `wb_user`;
CREATE TABLE `wb_user`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名称',
  `mobile` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `create_time` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '注册时间',
  `update_time` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `last_login_time` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最后登录时间',
  `status` tinyint(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态   1正常',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_mobile`(`mobile`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 31 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of wb_user
-- ----------------------------
INSERT INTO `wb_user` VALUES (30, '梁彬城', '18628329335', '313233343536d41d8cd98f00b204e9800998ecf8427e', 1588914013, 1588914013, 1590661029, 1);

SET FOREIGN_KEY_CHECKS = 1;
