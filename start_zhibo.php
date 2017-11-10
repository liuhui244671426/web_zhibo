<?php
//php在线直播示例代码
//命令：php start_zhibo.php

//设置路径
define('_ROOT_', dirname(__FILE__));
require_once _ROOT_.'/function.php';
//监听地址和端口
$server = new swoole_websocket_server("0.0.0.0", 8888);
//服务端接收连接事件
$server->on('open', function (swoole_websocket_server $server, $request) {
    if(!file_exists(_ROOT_.'/client/'.$request->fd.'.client')){
        @file_put_contents(_ROOT_.'/client/'.$request->fd.'.client',$request->fd);
    }
});
//服务端接收信息事件
$server->on('message', function (swoole_websocket_server $server, $frame) {
    //取fd, 群推
    foreach(notice(_ROOT_.'/client/') as $v){
            $server->push($v,$frame->data);
    }
});
//服务端接收关闭事件
$server->on('close', function ($ser, $fd) {
    @unlink(_ROOT_.'/client/'.$fd.'.client');
});
//服务开启
$server->start();