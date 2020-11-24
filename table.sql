CREATE TABLE `words`.`words`
(
    `id`         int(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `word`       varchar(255)    NOT NULL DEFAULT '' COMMENT '单词',
    `frequency`  int             NOT NULL COMMENT '出现频率',
    `createtime` datetime(0)     NULL,
    PRIMARY KEY (`id`)
);