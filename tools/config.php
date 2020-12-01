<?php


namespace tools;


class config
{
    protected static $config = null;

    /**
     * @param $name
     * @param $default
     * @return mixed
     */
    public static function get($name, $default = '')
    {
        if (is_null(self::$config)) {
            self::$config = require_once dirname(__DIR__) . '/env.php';
        }

        if (!isset(self::$config[$name])) {
            return $default;
        }

        return self::$config[$name];
    }
}