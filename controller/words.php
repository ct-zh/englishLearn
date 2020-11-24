<?php


namespace controller;


use tools\request;
use model\Words as model;

class words
{
    public static function save()
    {
        $words = request::clis();
        foreach ($words as $value) {
            $result = model::insert([
                'word' => $value,
                'frequency' => 0,
                'createtime' => date("Y-m-d", time())
            ]);
            var_dump($result);
        }
    }
}