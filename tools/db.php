<?php


namespace tools;

use PDO;
use PDOException;

class db
{
    /**
     * @var null|PDO
     */
    protected static $db = null;

    public static function getInstance()
    {
        if (is_null(self::$db)) {
            try {
                $pdo = new PDO(
                    config::get('dsn'),
                    config::get('username'),
                    config::get('password'),
                    config::get('options', null));
            } catch (PDOException $e) {
                echo json_encode([
                    'code' => -1,
                    'msg' => 'pdo connect error: ' . $e->getMessage()
                ]);
                die;
            }

            self::$db = $pdo;
        }
        return self::$db;
    }

}