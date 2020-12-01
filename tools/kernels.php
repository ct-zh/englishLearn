<?php


namespace tools;


use controller\words;
use tools\exception\error;

class kernels
{
    public static function run()
    {
        try {
            // 路由寻址
            $result = route::newInstance()->run();
        } catch (error $e) {
            exit(json_encode(['code' => -1, 'msg' => $e->getMessage()]));
        }

        if (!isset($result['code'])) {
            $result = ['code' => 0, 'data' => $result];
        }
        exit(json_encode($result));
    }

    public static function words()
    {
        words::save();
    }
}