<?php
namespace bybzmt\phpim;

class phpim
{
	const ACTION_SEND_MSG = 1;
	const ACTION_CLOSE_CONN = 2;
	const ACTION_BROADCAST = 3;
	const ACTION_ROOM_BROADCAST = 4;
	const ACTION_ROOM_ADD_CONN = 5;
	const ACTION_ROOM_DEL_CONN = 6;

	private $_actions = array();

	public function onCallback($on_connect, $on_msg, $on_disconnect)
	{
		$act = isset($_REQUEST['act']) ? $_REQUEST['act'] : '';

		$this->clearActions();

		switch ($act) {
		case 'connect' :
			$resp = $on_connect($this);
			break;

		case 'msg' :
			$id = isset($_REQUEST['id']) ? $_REQUEST['id'] : '';
			$msg = isset($_REQUEST['msg']) ? $_REQUEST['msg'] : '';

			$resp = $on_msg($this, $id, $msg);
			break;

		case 'disconnect' :
			$id = isset($_REQUEST['id']) ? $_REQUEST['id'] : '';
			$resp = $on_disconnect($this, $id);
			break;

		default:
			$resp = $this->setFail();
		}

		echo json_encode($resp);
	}

	public function doRequest($ip, $port)
	{
		$url = "http://{$ip}:{$port}/actions";

		$opts = array(
			'http'=>array(
				'method'=> 'POST',
				'header' => 'Content-type: application/json',
				'content' => json_encode($this->_actions),
				'ignore_errors' => true,
			)
		);

		$context = stream_context_create($opts);
		file_get_contents($url, false, $context);
	}

	public function setOk($id=null)
	{
		return array(
			'Ret' => 1,
			'Id' => $id,
			'Actions' => $this->_actions,
		);
	}

	public function setFail()
	{
		return array(
			'Ret' => 1,
			'Id' => null,
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
				'Conn' => $conn_id,
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
				'Conn' => $conn_id,
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
			),
		);

		return $this;
	}

	public function RoomBroadcast($room_id, $msg)
	{
		$this->_actions[] = array(
			'Type' => self::ACTION_ROOM_BROADCAST,
			'Point' => array(
				'Room' => $room_id,
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
				'Room' => $room_id,
				'Msg' => $msg,
			),
		);

		return $this;
	}

	public function RoomDelConn($room_id, $conn_id)
	{
		$this->_actions[] = array(
			'Type' => self::ACTION_ROOM_DEL_CONN,
			'Point' => array(
				'Room' => $room_id,
				'Msg' => $msg,
			),
		);

		return $this;
	}
}
