<?php

namespace model;

use tools\model;

class Words extends model
{
    protected static $table = 'words';

    public static function getWords($words)
    {
        array_walk($words, function (&$e) {
            $e = "'$e'";
        });

        $words = implode(", ", $words);
        $sql = "SELECT * FROM words where word in ({$words})";
        $query = self::db()->prepare($sql);

        $query->execute();
        return $query->fetchAll(self::db()::FETCH_ASSOC);
    }

    public static function incr($id)
    {
        $sql = "UPDATE words SET `frequency`=`frequency`+1 where id = '{$id}'";

        $query = self::db()->prepare($sql);
        $query->execute();
    }
}