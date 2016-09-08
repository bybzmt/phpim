phpim
========
这是一个迷你的im服务端。

服务端由go语言编写

附带一个php的api

这是一个最小做的im，服务端只会在有新连接，有新消息，有连接断开时
通过回调接口通知php，然后由php执行接下去的动作，go服务端本身不会
有动作。

api支持的动作包括

* 单播发送消息
* 连接加入房间
* 连接离开房间
* 房间广播消息
* 全服广播消息

phpclient
-------
这是php的客户端程序。
可以看demo.php里的代码。

composer安装
composer require "bybzmt/phpim"
