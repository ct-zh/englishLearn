<?php


namespace tools;


class request
{
    public static function get()
    {

    }

    public static function post()
    {

    }

    public static function cli($name)
    {

    }

    public static function clis()
    {
        $params = $_SERVER['argv'];
        array_shift($params);
        return $params;
    }
}