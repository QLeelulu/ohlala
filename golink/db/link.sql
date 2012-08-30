CREATE SCHEMA IF NOT EXISTS `link` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ; 
USE `link`; 

-- ----------------------------------------------------- 
-- Table `user` 用户表
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `user` ( 
  `id` BIGINT NOT NULL AUTO_INCREMENT , 
  `name` VARCHAR(100) NOT NULL , -- 用户名
  `email` VARCHAR(100) NOT NULL , -- email
  `email_lower` VARCHAR(100) NOT NULL , -- email小写，唯一键
  `pwd` CHAR(50) NOT NULL , -- 密码
  `user_pic` VARCHAR(1000) NOT NULL , -- 用户头像
  `description` VARCHAR(1000) NOT NULL , -- 自我介绍
  `reference_id` VARCHAR(1000) NOT NULL , -- 关联微博帐户id
  `reference_system` INT NOT NULL DEFAULT 0 , -- 微博平台类型
  `reference_token` VARCHAR(50) NOT NULL , -- 微博access token
  `reference_token_secret` VARCHAR(50) NOT NULL , -- 微博access token secret
  `link_count` INT(11) NOT NULL DEFAULT 0 , -- 分享的链接数量
  `friend_count` INT(11) NOT NULL DEFAULT 0 , -- 关注的数量
  `follower_count` INT(11) NOT NULL DEFAULT 0 , -- 粉丝的数量
  `topic_count` INT(11) NOT NULL DEFAULT 0 , -- 创建的话题的数量
  `ftopic_count` INT(11) NOT NULL DEFAULT 0 , -- 关注的话题的数量
  `create_time` datetime NOT NULL, -- 注册时间
  PRIMARY KEY (`id`) ,  
  INDEX `idx_reference_id` USING BTREE (`reference_id` ASC) , 
  INDEX `idx_name` USING BTREE (`name` ASC),
  UNIQUE INDEX `idx_email_lower` USING BTREE (`email_lower`),
  INDEX `idx_email_pwd` USING BTREE (`email_lower`,`pwd`) )
ENGINE = InnoDB, AUTO_INCREMENT = 10000;

-- -----------------------------------------------------
-- Table `user_follow` 用户跟随表
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `user_follow` (
  `user_id` BIGINT NOT NULL DEFAULT 0 , -- 跟随者的id
  `follow_id` BIGINT NOT NULL DEFAULT 0 ,-- 被跟随者的id
  `create_time` datetime NOT NULL, -- 跟随的时刻
  UNIQUE INDEX `idx_user_id` USING BTREE (`user_id`, `follow_id` ASC),
  INDEX `idx_follow_id` USING BTREE (`follow_id`, `user_id` ASC) )
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `topic_follow` 用户关注的话题
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `topic_follow` (
  `user_id` BIGINT NOT NULL DEFAULT 0 , -- 用户的id
  `topic_id` BIGINT NOT NULL DEFAULT 0 ,-- topic的id
  `create_time` datetime NOT NULL, -- 跟随的时刻
  UNIQUE INDEX `idx_user_id` USING BTREE (`user_id`, `topic_id` ASC),
  INDEX `idx_topic_id` USING BTREE (`topic_id`) )
ENGINE = InnoDB;

-- ----------------------------------------------------- 
-- Table `link` 分享链接表
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `link` ( 
  `id` BIGINT NOT NULL AUTO_INCREMENT , 
  `user_id` BIGINT NOT NULL , -- 用户id
  `title` VARCHAR(200) NOT NULL , -- 链接标题
  `context` VARCHAR(500) NOT NULL , -- 链接内容（链接、文本内容）
  `context_type` INT NOT NULL DEFAULT 0 , -- 内容类型（链接、文本）
  `topics` VARCHAR(500) NOT NULL , -- 标签已分号隔开
  `create_time` DATETIME NOT NULL , -- 创建时间
  `vote_up` BIGINT NOT NULL DEFAULT 0 ,-- 顶的数量
  `vote_down` BIGINT NOT NULL DEFAULT 0 , -- 踩的数量
  `reddit_score` DECIMAL(28,10) NOT NULL , -- 链接得分
  `view_count` INT(11) unsigned NOT NULL DEFAULT 0 , -- 链接查看次数
  `comment_count` INT(11) unsigned NOT NULL DEFAULT 0 , -- 链接评论数
   -- `comment_reddit_score` DECIMAL(28,10) NOT NULL ,
  PRIMARY KEY (`id` DESC) , 
  INDEX `idx_title` USING BTREE (`title` ASC),
  INDEX `idx_user_id` USING BTREE (`user_id` ASC),
  INDEX `idx_create_time` USING BTREE (`create_time` DESC)
  )
ENGINE = InnoDB, AUTO_INCREMENT = 10001; 

-- ----------------------------------------------------- 
-- Table `link_host` 链接的host关系表
-- ----------------------------------------------------- 
CREATE TABLE `link_host` (
  `host_id` bigint(20) unsigned NOT NULL,
  `link_id` bigint(20) unsigned NOT NULL,
  INDEX `idx_host_id` USING BTREE (`host_id` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------------------------------- 
-- Table `topic` 话题表
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `topic` ( 
  `id` BIGINT NOT NULL AUTO_INCREMENT , 
  `name` VARCHAR(50) NOT NULL , -- 话题名称
  `name_lower` VARCHAR(50) NOT NULL , -- 话题名小写，唯一索引
  `description` VARCHAR(250) NULL , -- 话题的描述
  `pic` VARCHAR(100) NULL , -- 话题的图片
  `click_count` BIGINT NOT NULL DEFAULT 0 , -- 话题点击次数
  `follower_count` BIGINT NOT NULL DEFAULT 0 , -- 话题的关注者数量
  `link_count` BIGINT NOT NULL DEFAULT 0 , -- 添加到该话题的链接数量
  PRIMARY KEY (`id` DESC),
  UNIQUE INDEX `idx_name_lower` USING BTREE (`name_lower`) ) 
ENGINE = InnoDB, AUTO_INCREMENT = 10001;

-- ----------------------------------------------------- 
-- Table `topic_link` 标签与链接表关联
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `topic_link` ( 
  `topic_id` BIGINT NOT NULL DEFAULT 0 , -- 标签id
  `link_id` BIGINT NOT NULL DEFAULT 0 , -- 链接id
  -- INDEX `idx_topic_id` USING BTREE (`topic_id` ASC) 
  UNIQUE INDEX `idx_topic_link` USING BTREE (`topic_id`,`link_id`), 
  INDEX `idx_link_id` USING BTREE (`link_id` DESC)
  ) 
ENGINE = InnoDB; 


-- ----------------------------------------------------- 
-- Table `u_Comment` 评论表
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `comment` ( 
  `id` BIGINT NOT NULL AUTO_INCREMENT , 
  `link_id` BIGINT NOT NULL DEFAULT 0 , -- link的id
  `user_id` BIGINT NOT NULL DEFAULT 0 , -- user的id
  `parent_path` VARCHAR(2000) NOT NULL , -- 父节点id路径,根节点为空字符
  `children_count` INT NOT NULL DEFAULT 0 , -- 子节点个数(当前一级)
  `top_parent_id` BIGINT NOT NULL DEFAULT 0 , -- 根节点id
  `parent_id` BIGINT NOT NULL DEFAULT 0 , -- 父节点id
  `deep` tinyint(4) unsigned NOT NULL DEFAULT '0', -- 节点深度
  `status` INT NOT NULL DEFAULT 1 , -- 评论状态：1代表正常、2代表删除
  `content` VARCHAR(1000) NOT NULL , -- 评论内容
  `create_time` DATETIME NOT NULL , -- 评论时间
  `vote_up` BIGINT NOT NULL DEFAULT 0 ,-- 支持加数
  `vote_down` BIGINT NOT NULL DEFAULT 0 , -- 支持减数
  `reddit_score` DECIMAL(28,10) NOT NULL , -- 根节点评论得分
  -- `children_reddit_score` DECIMAL(28,10) NOT NULL , -- 子节点评论得分总和，只有根节点才有值，子节点该字段值为0
  PRIMARY KEY (`id` DESC) , 
  INDEX `idx_link_id` USING BTREE (`link_id` ASC), 
  INDEX `idx_top_parent_id` USING BTREE (`top_parent_id`,`parent_id` ASC) ) 
ENGINE = InnoDB, AUTO_INCREMENT = 10000;

-- ----------------------------------------------------- 
-- Table `u_LinkSupportRecord` 用户支持与链接表的关联表
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `link_support_record` ( 
  `link_id` BIGINT NOT NULL DEFAULT 0 , -- 链接表的id
  `user_id` BIGINT NOT NULL DEFAULT 0 , -- 用户id
  `score` INT NOT NULL DEFAULT 0 , -- 得分（正负一）
  INDEX `idx_link_id` USING BTREE (`link_id`,`user_id` ASC)
  --  , INDEX `IDX_FUserID` USING BTREE (`FUserID` ASC) 
  ) 
ENGINE = InnoDB; 


-- ----------------------------------------------------- 
-- Table `u_CommentSupportRecord` 用户支持与评论表的关联表
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `comment_support_record` ( 
  `comment_id` BIGINT NOT NULL DEFAULT 0 , -- 评论id
  `user_id` BIGINT NOT NULL DEFAULT 0 , -- 用户id
  `score` INT NOT NULL DEFAULT 0 , -- 得分（正负一）
  INDEX `idx_comment_id` (`comment_id`,`user_id` ASC)
  -- , INDEX `IDX_FUserID` USING BTREE (`FUserID` ASC) 
  ) 
ENGINE = InnoDB; 


-- ----------------------------------------------------- 
-- Table `tui_link_for_topic_later` 从某个话题去浏览最新链接的推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_topic_later` (
  `topic_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_topic_link` USING BTREE (`topic_id`,`link_id`)
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `tui_link_for_topic_top` 从某个话题去浏览热门链接的推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_topic_top` (
  `topic_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `create_time` datetime NOT NULL,
  `reddit_score` DECIMAL(28,10) NOT NULL , -- 热门的排序
  UNIQUE KEY `idx_topic_link` USING BTREE (`topic_id`,`link_id`)
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `tui_link_for_topic_hot` 从某个话题去浏览热议链接的推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_topic_hot` (
  `topic_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `time_type` int NOT NULL DEFAULT 0 ,-- 投票时间范围: 1:全部时间；2:这个小时；3:今天；4:这周；5:这个月；6:今年
  `vote_abs_score` int NOT NULL , -- 热议的排序,|up - down| 趋向于0代表热议
  `vote_add_score` int NOT NULL , -- 热议的排序,up + down 越大代表热议
  `create_time` datetime NOT NULL,
  UNIQUE KEY `idx_topic_link` USING BTREE (`topic_id`,`link_id`,`time_type`)
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `tui_link_for_topic_vote` 从某个话题去浏览投票高链接的推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_topic_vote` (
  `topic_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `time_type` int NOT NULL DEFAULT 0 ,-- 投票时间范围: 1:全部时间；2:这个小时；3:今天；4:这周；5:这个月；6:今年
  `vote` int NOT NULL DEFAULT 0 ,-- up - down 越大越靠前
  `create_time` datetime NOT NULL,
  UNIQUE KEY `idx_topic_link` USING BTREE (`topic_id`,`link_id`,`time_type`)
) ENGINE=InnoDB;


-- ----------------------------------------------------- 
-- Table `tui_link_for_home` 首页的推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_home` (
  `link_id` bigint(20) NOT NULL,
  `data_type` int NOT NULL, -- 2:热门; 3:热议[3:全部时间；10:这个小时；11:今天；12:这周；13:这个月；14:今年]; [投票时间范围: 4:全部时间；5:这个小时；6:今天；7:这周；8:这个月；9:今年]
  `score` DECIMAL(28,10) NOT NULL , -- 各种排序的得分
  `vote_add_score` int NOT NULL DEFAULT 0, -- 热议的排序,up + down 越大代表热议
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_topic_link` USING BTREE (`data_type`, `link_id`)
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `tui_link_for_handle` 链接处理队列表 
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_handle` ( 
`link_id` bigint(20) NOT NULL DEFAULT 0, 
`create_time` datetime NOT NULL, -- 链接的发布时间 
`user_id` bigint(20) NOT NULL DEFAULT 0, -- 发布者的id，如果是投票就不需要 
`insert_time` datetime NOT NULL, -- 记录插入的时间 
`data_type` int NOT NULL, -- 1:新增; 2:投票; 
PRIMARY KEY (`link_id` DESC), 
INDEX `idx_insert_time` USING BTREE (`insert_time` DESC) 
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `tui_link_for_delete` 需要删除的益处数据 
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_delete` ( 
`id` bigint(20) NOT NULL , 
`time_type` int NOT NULL DEFAULT 0, -- 数据类型 
`del_count` bigint(20) NOT NULL 
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `tui_link_temporary_delete` 需要删除的临时表 
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_temporary_delete` ( 
`id` bigint(20) NOT NULL
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `link_for_user` 用户及话题的链接推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `link_for_user` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `link_for_user_0` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_1` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_2` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_3` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_4` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_5` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_6` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_7` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_8` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_9` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_10` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_11` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_12` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_13` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_14` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_15` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_16` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_17` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_18` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_19` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_20` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_21` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_22` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_23` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS `link_for_user_24` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE INDEX `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;


