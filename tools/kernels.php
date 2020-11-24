<?php


namespace tools;


use controller\words;

class kernels
{
    public static function words()
    {
        words::save();
    }
}