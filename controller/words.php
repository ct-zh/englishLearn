<?php


namespace controller;


use tools\request;
use model\Words as model;

class words
{
    public static function save()
    {
        $newWords = request::clis();

        $exist = model::getWords($newWords);
        $exist = array_column($exist, null, "word");

        foreach ($newWords as $value) {
            // todo: 单词转换， 复数形式/过去式 等等
            if (isset($exist[$value])) {
                model::incr($exist[$value]["id"]);
            } else {
                model::insert([
                    'word' => $value,
                    'frequency' => 0,
                    'createtime' => date("Y-m-d H:i:s", time())
                ]);
            }
        }

        echo json_encode(['code' => 0, 'msg' => 'ok']) . PHP_EOL;
        die;
    }
}