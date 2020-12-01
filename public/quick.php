<?php
/**
 * 单文件快速开发
 */

if (PHP_SAPI == 'cli') {
    exit("暂不支持命令行" . PHP_EOL);
}

$loadFile = dirname(__DIR__) . '/tools/load.php';
if (!file_exists($loadFile)) {
    exit("配置文件加载出错");
}
$res = include $loadFile;
tools\load::autoload();

// 路由
$quick = new QuickDo();
switch ($_SERVER['QUERY_STRING']) {
    case "d":
        $result = $quick->daily();
        break;
    default:
        $result = ['code' => -1, 'msg' => "not found"];
}
exit(json_encode($result));


class QuickDo
{
    // 保存一个日常任务流程
    public function daily()
    {
        $startTime = $_POST["starttime"];
        $endTime = $_POST["endtime"];
        $content = trim($_POST["content"]);
        if (strlen($content) > 1000) {
            return ['code' => -1, "msg" => "内容不能超过1000字节"];
        }

        $startTime = date('Y-m-d H:i:s', strtotime($startTime));
        $data = [
            'content' => $content,
            'starttime' => $startTime,
            'type' => \enum\week::TYPE_DAILY,
        ];
        if (empty($endTime)) {
            $data['status'] = \enum\week::STATUS_DOING;
        } else {
            $endTime = date('Y-m-d H:i:s', strtotime($endTime));
            $data['endtime'] = $endTime;
            $data['status'] = \enum\week::STATUS_DONE;
        }

        $result = \model\week::insert($data);
        if (!$result) {
            return ['code' => -1, 'msg' => "insert error"];
        }

        return ['code' => 0, 'msg' => 'success', 'data' => $result];
    }
}