<?php

namespace enum;

class week
{
    const STATUS_DOING = 1;
    const STATUS_DONE = 2;
    const STATUS_DELETE = 3;

    const STATUS_MAP = [
        self::STATUS_DOING => "正在进行",
        self::STATUS_DONE => "已完成",
        self::STATUS_DELETE => '被删除',
    ];

    CONST TYPE_DAILY = 1;
    const TYPE_WEEKLY = 2;

    const TYPE_MAP = [
        self::TYPE_DAILY => "日常任务",
        self::TYPE_WEEKLY => "周任务",
    ];

}