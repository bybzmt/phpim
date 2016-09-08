<?php
namespace bybzmt\phpim;

/**
 * im服务封装
 */
class phpim
{
	const ACTION_SEND_MSG = 1;
	const ACTION_CLOSE_CONN = 2;
	const ACTION_BROADCAST = 3;
	const ACTION_ROOM_BROADCAST = 4;
	const ACTION_ROOM_ADD_CONN = 5;
	const ACTION_ROOM_DEL_CONN = 6;

	static public $host;

	private $_actions = array();

	/**
	 * callback时动作
	 *
	 * @param callable on_connect($im) 在有新连接时回调
	 * @param callable on_msg($im, $conn_id, $msg) 新消息时回调
	 * @param callable on_disconnect($im, $conn_id) 连接断开时回调
	 */
	public function onCallback(callable $on_connect, callable $on_msg, callable $on_disconnect)
	{
		$act = isset($_REQUEST['act']) ? $_REQUEST['act'] : '';

		$this->clearActions();

		switch ($act) {
		case 'connect' :
			$resp = call_user_func($on_connect, $this);
			break;

		case 'msg' :
			$id = isset($_REQUEST['id']) ? $_REQUEST['id'] : '';
			$msg = isset($_REQUEST['msg']) ? $_REQUEST['msg'] : '';

			$resp = call_user_func($on_msg, $this, $id, $msg);
			break;

		case 'disconnect' :
			$id = isset($_REQUEST['id']) ? $_REQUEST['id'] : '';
			$resp = call_user_func($on_disconnect, $this, $id);
			break;

		default:
			$resp = $this->setFail();
		}

		echo json_encode($resp);
	}

	public function doRequest($host=null)
	{
		if (!$host) {
			$host = self::$host;
		}
		$url = "http://{$host}/actions";

		$opts = array(
			'http'=>array(
				'method'=> 'POST',
				'header' => implode("\r\n", array(
					'Content-type: application/json',
					"Connection: close",
				)),
				'content' => json_encode($this->_actions),
				'ignore_errors' => true,
			)
		);

		$context = stream_context_create($opts);
		$out = file_get_contents($url, false, $context);
		return $out;
	}

	public function setOk($id=null)
	{
		return array(
			'Ret' => 0,
			'Id' => "$id",
			'Actions' => $this->_actions,
		);
	}

	public function setFail()
	{
		return array(
			'Ret' => 1,
			'Id' => '',
			'Actions' => $this->_actions,
		);
	}

	public function clearActions()
	{
		$this->_actions = array();
	}

	public function sendMsg($conn_id, $msg)
	{
		$this->_actions[] = array(
			'Type' => self::ACTION_SEND_MSG,
			'Point' => array(
				'Conn' => "$conn_id",
				'Msg' => $msg,
			),
		);

		return $this;
	}

	public function CloseConn($conn_id)
	{
		$this->_actions[] = array(
			'Type' => self::ACTION_CLOSE_CONN,
			'Point' => array(
				'Conn' => "$conn_id",
			),
		);

		return $this;
	}

	public function Broadcast($msg)
	{
		$this->_actions[] = array(
			'Type' => self::ACTION_BROADCAST,
			'Point' => array(
				'Msg' => $msg,
			)
		);

		return $this;
	}

	public function RoomBroadcast($room_id, $msg)
	{
		$this->_actions[] = array(
			'Type' => self::ACTION_ROOM_BROADCAST,
			'Point' => array(
				'Room' => "$room_id",
				'Msg' => $msg,
			),
		);

		return $this;
	}

	public function RoomAddConn($room_id, $conn_id)
	{
		$this->_actions[] = array(
			'Type' => self::ACTION_ROOM_ADD_CONN,
			'Point' => array(
				'Room' => "$room_id",
				'Conn' => $conn_id,
			),
		);

		return $this;
	}

	public function RoomDelConn($room_id, $conn_id)
	{
		$this->_actions[] = array(
			'Type' => self::ACTION_ROOM_DEL_CONN,
			'Point' => array(
				'Room' => "$room_id",
				'Conn' => $conn_id,
			),
		);

		return $this;
	}
}
