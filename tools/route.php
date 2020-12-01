<?php


namespace tools;

use tools\exception\error;
use tools\extra\singleton;

class route
{
    use singleton;

    protected static $route = [];

    /**
     * @throws error
     */
    public function run()
    {
        // 加载路由配置
        $this->loadConfig();

        // 问题： 参数只能存放在form data 里面， get请求不能获取参数,
        // 现在先只用post方法把
        $classPath = self::$route[$_SERVER["REQUEST_METHOD"]][$_SERVER["QUERY_STRING"]] ?? null;
        if ($classPath == null) {
            throw new error("404 NOT FOUND");
        }

        // 解析路由
        [$class, $method] = explode("@", $classPath);
        $class = "controller\\" . $class;
        try {
            $reflectionClass = new \ReflectionClass($class);
            $reflectMethod = $reflectionClass->getMethod($method);
            $result = $reflectMethod->invoke(new $class());
        } catch (\ReflectionException $e) {
            throw new error("route error: " . $e->getMessage());
        }

        return $result;
    }

    protected function loadConfig()
    {
        $path = dirname(__DIR__) . "/config/route.php";
        if (file_exists($path)) {
            require_once $path;
        }
    }

    public static function addRoute($method, $httpPath, $classPath)
    {
        // check
        if (!in_array($method, ["post", "put", "delete", "get"])) {
            $method = "get";
        }
        self::$route[strtoupper($method)][$httpPath] = $classPath;
    }
}