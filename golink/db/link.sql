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
  `create_time` datetime NOT NULL, -- 注册时间
  PRIMARY KEY (`id`) , 
  INDEX `idx_reference_id` USING BTREE (`reference_id` ASC) , 
  INDEX `idx_name` USING BTREE (`name` ASC),
  UNIQUE KEY `idx_email_lower` USING BTREE (`email_lower`),
  INDEX `idx_email_pwd` USING BTREE (`email_lower`,`pwd`) )
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `user_follow` 用户跟随表
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `user_follow` (
  `user_id` BIGINT NOT NULL DEFAULT 0 , -- 跟随者的id
  `follow_id` BIGINT NOT NULL DEFAULT 0 ,-- 被跟随者的id
  `create_time` datetime NOT NULL, -- 跟随的时刻
  INDEX `idx_user_id` USING BTREE (`user_id`, `follow_id` ASC),
  INDEX `idx_follow_id` USING BTREE (`follow_id`, `user_id` ASC) )
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `topic_follow` 用户关注的话题
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `topic_follow` (
  `user_id` BIGINT NOT NULL DEFAULT 0 , -- 用户的id
  `topic_id` BIGINT NOT NULL DEFAULT 0 ,-- topic的id
  `create_time` datetime NOT NULL, -- 跟随的时刻
  INDEX `idx_user_id` USING BTREE (`user_id`, `topic_id` ASC),
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
   `comment_reddit_score` DECIMAL(28,10) NOT NULL ,
  PRIMARY KEY (`id` DESC) , 
  INDEX `idx_title` USING BTREE (`title` ASC),
  INDEX `idx_create_time` USING BTREE (`create_time` DESC)
  )
ENGINE = InnoDB; 


-- ----------------------------------------------------- 
-- Table `topic` 话题表
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `topic` ( 
  `id` BIGINT NOT NULL AUTO_INCREMENT , 
  `name` VARCHAR(50) NOT NULL , -- 话题名称
  `name_lower` VARCHAR(50) NOT NULL , -- 话题名小写，唯一索引
  `desc` VARCHAR(250) NULL , -- 话题的描述
  `pic` VARCHAR(100) NULL , -- 话题的图片
  `clicks` BIGINT NOT NULL DEFAULT 0 , -- 话题点击次数
  `followers` BIGINT NOT NULL DEFAULT 0 , -- 话题的关注者数量
  `links` BIGINT NOT NULL DEFAULT 0 , -- 添加到该话题的链接数量
  PRIMARY KEY (`id` DESC),
  UNIQUE KEY `idx_name_lower` USING BTREE (`name_lower`) ) 
ENGINE = InnoDB;

-- ----------------------------------------------------- 
-- Table `topic_link` 标签与链接表关联
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `topic_link` ( 
  `topic_id` BIGINT NOT NULL DEFAULT 0 , -- 标签id
  `link_id` BIGINT NOT NULL DEFAULT 0 , -- 链接id
  -- INDEX `idx_topic_id` USING BTREE (`topic_id` ASC) 
  UNIQUE KEY `idx_topic_link` USING BTREE (`topic_id`,`link_id`), 
  INDEX `idx_link_id` USING BTREE (`link_id` DESC)
  ) 
ENGINE = InnoDB; 


-- ----------------------------------------------------- 
-- Table `u_Comment` 评论表
-- ----------------------------------------------------- 
CREATE  TABLE IF NOT EXISTS `comment` ( 
  `id` BIGINT NOT NULL AUTO_INCREMENT , 
  `link_id` BIGINT NOT NULL DEFAULT 0 , -- link的id
  `top_parent_id` BIGINT NOT NULL DEFAULT 0 , -- 根节点id
  `parent_id` BIGINT NOT NULL DEFAULT 0 , -- 父节点id
  `status` INT NOT NULL DEFAULT 0 , -- 评论状态：1代表正常、2代表删除
  `content` VARCHAR(1000) NOT NULL , -- 评论内容
  `create_time` DATETIME NOT NULL , -- 评论时间
  `vote_up` BIGINT NOT NULL DEFAULT 0 ,-- 支持加数
  `vote_down` BIGINT NOT NULL DEFAULT 0 , -- 支持减数
  `reddit_score` DECIMAL(28,10) NOT NULL , -- 根节点评论得分
  `children_reddit_score` DECIMAL(28,10) NOT NULL , -- 子节点评论得分总和，只有根节点才有值，子节点该字段值为0
  PRIMARY KEY (`id` DESC) , 
  INDEX `idx_link_id` USING BTREE (`link_id` ASC), 
  INDEX `idx_top_parent_id` USING BTREE (`top_parent_id`,`parent_id` ASC) ) 
ENGINE = InnoDB;

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
-- Table `link_for_user` 用户及话题的链接推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `link_for_user` (
  `user_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  `user_count` INT NOT NULL,
  `topic_count` INT NOT NULL,
  -- `data_type` int NOT NULL, -- 1:关注者的推送；2:话题的推送；3:关注者与话题的推送 [控制1和3的记录和<=1w; 2和3的记录一样控制]
  `create_time` datetime NOT NULL,
  UNIQUE KEY `idx_user_link` USING BTREE (`user_id`,`link_id`)
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `tui_link_for_topic` 从某个话题去浏览链接的推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_topic` (
  `topic_id` bigint(20) NOT NULL,
  `link_id` bigint(20) NOT NULL,
  UNIQUE KEY `idx_topic_link` USING BTREE (`topic_id`,`link_id`)
) ENGINE=InnoDB;

-- ----------------------------------------------------- 
-- Table `tui_link_for_home` 首页的推送表
-- ----------------------------------------------------- 
CREATE TABLE IF NOT EXISTS `tui_link_for_home` (
  `link_id` bigint(20) NOT NULL,
  `data_type` int NOT NULL, -- 1:最新; 2:热门; 3:热议; 4:得分
  `vote_score` BIGINT NOT NULL DEFAULT 0 ,-- 投票数之和
  `reddit_score` DECIMAL(28,10) NOT NULL , -- 热门的排序
  `comment_reddit_score` DECIMAL(28,10) NOT NULL , -- 热议的排序
  UNIQUE KEY `idx_topic_link` USING BTREE (`data_type`, `link_id`)
) ENGINE=InnoDB;




