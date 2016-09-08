<?php
require __DIR__ . '/phpim.php';

use bybzmt\phpim;

#回调接口
#im服务会在新连接，新消息，连接断开时回调php的接口
#
#当有新连接时php需要决定是否充许连接，充许连接时需要指定一个
#唯一的id,之后对这个连接的操作都需要通过这个唯一的id进行
#如果id重复，之前的连接将会丢失。
#
#当有新消息时,php会收到发前消息的连接的id和消息内容，php需要
#响应接下去的动作，im本身不会做任何的动作,php就响应就相当于
#忽略消息
#
#当连接断开时，php会收到断开连接的id
#
$on_connect = function($im) {
	$conn_id = 'conn_id';
	$room_id = 'room_id';

	//可选动作
	$im->sendMsg($conn_id, '新消息');
	$im->RoomAddConn($room_id, $conn_id);
	$im->RoomBroadcast($room_id, '房间内消息');

	//成功, 充许连接，连接成功后会执行上面添加的命令
	return $im->setOk($conn_id);
	//失败，阻止对方连接上来
	//return $im->setFail();
};
$on_msg = function($im, $id, $msg) {
	//把消息向所有人广播
	$im->Broadcast($msg);
	return $im->setOk();
};
$on_disconnect = function($im) {
	//这时候其实只是通知连接断了，并不能对这个连接做什么。
	//当然可以对其它的连接发消息
	return $im->setOk();
};

$im = new phpim();
$im->onCallback($on_connect, $on_msg, $on_disconnect);
