<?php
/**
 * 单例模式
 */

namespace tools\extra;


trait singleton
{
    public static function newInstance()
    {
        static $self = null;
        if ($self == null) {
            $self = new static();
        }
        return $self;
    }
}