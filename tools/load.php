<?php


namespace tools;


class load
{
    protected static $classAlias = [

    ];

    public static function autoload()
    {
        spl_autoload_register('tools\\load::defaultLoad');
    }

    public static function defaultLoad($class)
    {
        if (isset(self::$classAlias[$class])) {
            return class_alias(self::$classAlias[$class], $class);
        }

        if ($file = self::findFile($class)) {
            __include_file($file);
            return true;
        }
    }

    protected static function findFile($class)
    {
        $class = str_replace('\\', '/', $class);
        $file = dirname(__DIR__) . "/$class.php";
        if (file_exists($file)) {
            return $file;
        }
        return false;
    }
}

function __include_file($file)
{
    return include_once $file;
}