SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
  `status` tinyint(3) unsigned NOT NULL COMMENT '状态位(1:成功,2:失败,3:待重试,4:重试中,9:被忽略)',
  `duration` decimal(16,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '投递耗时',
  `retry` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '投递次数',
  `task_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订阅任务ID',
  `payload_message_id` varchar(32) DEFAULT NULL COMMENT '主题消息ID',
  `message_dequeue` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '第几次出列时落库',
  `message_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '发布时间',
  `message_id` varchar(32) NOT NULL COMMENT '消息ID',
  `message_body` text NOT NULL COMMENT '消息正文',
  `response_body` text COMMENT '消息投递结果',
  `gmt_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_task_message` (`task_id`,`message_id`) USING BTREE,
  KEY `idx_topic_message` (`payload_message_id`),
  KEY `idx_status` (`status`,`task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8 COMMENT='消费记录';

-- ----------------------------
-- Table structure for payload
-- ----------------------------
DROP TABLE IF EXISTS `payload`;
CREATE TABLE `payload` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '状态位(1:成功,2:失败,3:待重试,4:重试中)',
  `duration` decimal(16,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '发布耗时',
  `retry` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发布次数',
  `message_task_id` int(10) NOT NULL DEFAULT '0' COMMENT '消费记录表的task_id字段值',
  `message_message_id` varchar(32) NOT NULL COMMENT '消费记录表的message_id字段值',
  `hash` char(32) NOT NULL COMMENT '发布哈希',
  `offset` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发布批次',
  `registry_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '注册关系表的id字段值',
  `message_id` varchar(32) DEFAULT NULL COMMENT 'MQ服务器返回的消息ID',
  `message_body` text NOT NULL COMMENT 'MQ消息内容',
  `response_body` text COMMENT 'MQ发布结果',
  `gmt_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_hash` (`hash`,`offset`) USING BTREE,
  KEY `idx_message_id` (`message_id`) USING BTREE,
  KEY `idx_topic_tag` (`registry_id`),
  KEY `idx_status` (`status`,`registry_id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8 COMMENT='生产记录';

-- ----------------------------
-- Table structure for registry
-- ----------------------------
DROP TABLE IF EXISTS `registry`;
CREATE TABLE `registry` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
  `topic_name` varchar(32) NOT NULL COMMENT '主题名',
  `topic_tag` varchar(64) NOT NULL COMMENT '主题标签',
  `filter_tag` varchar(16) DEFAULT NULL COMMENT '过滤标签',
  `gmt_created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_topic_pair` (`topic_name`,`topic_tag`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8 COMMENT='注册关系';

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '状态位(0:停用, 1:启用)',
  `title` varchar(128) NOT NULL COMMENT '标题',
  `remark` text COMMENT '描述',
  `parallels` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '最大并行数(单个节点最多开启消费者数)',
  `concurrency` tinyint(3) unsigned NOT NULL DEFAULT '32' COMMENT '最大并发数(单个消费者最多允许多少条消息同时处于投递中)',
  `max_retry` tinyint(3) unsigned NOT NULL DEFAULT '3' COMMENT '最大重试数(投递失败的消息, 最多允许重试次数)',
  `delay_seconds` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '消息发布后, 延时多久(秒)再允许消费',
  `broadcasting` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否广播(支持: Rocketmq)',
  `registry_id` int(10) unsigned NOT NULL COMMENT '注册关系ID',
  `handler` varchar(255) NOT NULL COMMENT '订阅回调地址',
  `handler_timeout` tinyint(3) NOT NULL DEFAULT '10' COMMENT '订阅回调超时(单位: 秒)',
  `handler_method` varchar(16) DEFAULT NULL COMMENT '订阅回调方式',
  `handler_condition` varchar(255) DEFAULT NULL COMMENT '条件过滤',
  `handler_response_type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '投递结果(0:JSON.ERRNO=0,1:HTML.CODE=200)',
  `handler_ignore_codes` varchar(255) DEFAULT NULL COMMENT '忽略状态码',
  `failed` varchar(255) DEFAULT NULL COMMENT '失败通知地址',
  `failed_timeout` tinyint(3) NOT NULL DEFAULT '10' COMMENT '失败通知超时(单位:秒)',
  `failed_method` varchar(16) DEFAULT NULL COMMENT '失败通知方式',
  `failed_condition` varchar(255) DEFAULT NULL,
  `failed_response_type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '失败回调结果类型(0:JSON.ERRNO=0,1:HTML.CODE=200)',
  `failed_ignore_codes` varchar(255) DEFAULT NULL COMMENT '失败忽略状态码',
  `succeed` varchar(255) DEFAULT NULL COMMENT '成功通知地址',
  `succeed_timeout` tinyint(3) NOT NULL DEFAULT '10' COMMENT '成功通知超时(单位:秒)',
  `succeed_method` varchar(16) DEFAULT NULL COMMENT '成功通知方式',
  `succeed_condition` varchar(255) DEFAULT NULL,
  `succeed_response_type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '成功回调结果类型(0:JSON.ERRNO=0,1:HTML.CODE=200)',
  `succeed_ignore_codes` varchar(255) DEFAULT NULL COMMENT '成功忽略状态码',
  `gmt_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `gmt_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uni_registry_handler` (`registry_id`,`handler`) USING BTREE,
  KEY `idx_status_updated` (`status`,`gmt_updated`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='订阅任务';

SET FOREIGN_KEY_CHECKS = 1;
