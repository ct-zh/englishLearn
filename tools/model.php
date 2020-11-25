<?php


namespace tools;


class model
{
    protected static $table;

    public static function insert($data, $option = '')
    {
        $fields = [];
        $bindVals = [];
        $bindKeys = [];
        foreach ($data as $key => $value) {
            $bindKeys[] = ":$key";
            $fields[] = "`{$key}`";
            $bindVals[":$key"] = $value;
        }

        $query = sprintf("INSERT %s INTO `%s` (%s) VALUES (%s)",
            $option, static::$table,
            implode(', ', $fields),
            implode(', ', $bindKeys));

        try {
            $sth = db::getInstance()->prepare($query);
            if (!$sth->execute($bindVals)) {
                return false;
            }
            $count = $sth->rowCount();
        } catch (\PDOException $e) {
            echo json_encode([
                'code' => -1,
                'msg' => 'sql error: ' . $e->getMessage()
            ]);
            die;
        }

        if ($count > 0) {
            return db::getInstance()->lastInsertId();
        }
        return 0;
    }

    public static function db()
    {
        return db::getInstance();
    }
}