var socket;
var startHight
var contract_hash = "6b5042927197210604d90935af76c1bb67807fff"
var reverse = function( str ){
   var stack = [];//生成一个栈
   for(var len = str.length,i=len;i>=0;i-- ){
       stack.push(str[i]);
    }
    return stack;
};

function str2ab(str) {
  var buf = new ArrayBuffer(str.length * 2); // 每个字符占用2个字节
  var bufView = new Uint16Array(buf);
  for (var i = 0, strLen = str.length; i < strLen; i++) {
    bufView[i] = str.charCodeAt(i);
  }
  return buf;
}
function hexToString(trimedStr){
　　var rawStr = trimedStr.substr(0,2).toLowerCase() === "0x" ? trimedStr.substr(2) : trimedStr;
　　var len = rawStr.length;
　　if(len % 2 !== 0) {
　　　　alert("Illegal Format ASCII Code!");
　　　　return "";
　　}
　　var curCharCode;
　　var resultStr = [];
　　for(var i = 0; i < len;i = i + 2) {
　　　　curCharCode = parseInt(rawStr.substr(i, 2), 16);
　　　　resultStr.push(String.fromCharCode(curCharCode));
　　}
　　return resultStr.join("");
}
function stringToHex(str){
	var val="";
	for(var i = 0; i < str.length; i++){
		if(val == "")
			val = str.charCodeAt(i).toString(16);
		else
			val += str.charCodeAt(i).toString(16);
	}
	return val;
　}
function GetEventsByHash(){
	req = {
    		"Action": "getsmartcodeeventbyhash",
    		"Version": "1.0.0",
    		"Id":125, //optional
    		"Hash": "755a99b57d999548594b0433bae86f1d0afa32566171888cd1ae14dc654c1c95",
	};
	socket.send(JSON.stringify(req));
}


function GetOrder(ID){
	req = {
    		"Action": "getstorage",
    		"Version": "1.0.0",
    		"Id":ID, //optional
    		"Hash": contract_hash,
    		"Key" : stringToHex("__ORDER___"+ID.toString(16)),
	};
	socket.send(JSON.stringify(req));
}
function GetOrders(ID){
	var id_str = ID.toString(16);
	if(id_str.length / 2 != 0)
		id_str = "0" + id_str;
	req = {
    		"Action": "getstorage",
    		"Version": "1.0.0",
    		"Id":1, //optional
    		"Hash": contract_hash,
    		"Key" : stringToHex("__ORDER___")+id_str,
	};
        //alert(stringToHex("__ORDER___")+id_str);
	socket.send(JSON.stringify(req));
}

$(document).ready(function(){
  $("#myCarousel").carousel();

    $.get('/api?req_type=query_order_rank&order_type=_BUY___List_Tail_Order___ONG_ONT_', function(data) {
      var arr = JSON.parse(data);
        console.log(arr);
      	for(var i in arr){
                $("#buy_table").append("<tr><td style='color:#CC0033'>"+arr[i].Price/100000000+"</td><td style='color:#CC0033'>"+(arr[i].Price*arr[i].UnAmount)/100000000+"</td><td style='color:#CC0033'>"+arr[i].UnAmount+"</td></tr>");
	}	
    });
    $.get('/api?req_type=query_order_rank&order_type=_BUY___List_Tail_Order___ONG_ONT_', function(data) {
      var arr = JSON.parse(data);
        console.log(arr);
      	for(var i in arr){
                $("#sell_table").append("<tr><td style='color:#669933'>"+arr[i].Price/100000000+"</td><td style='color:#669933'>"+(arr[i].Price*arr[i].UnAmount)/100000000+"</td><td style='color:#669933'>"+arr[i].UnAmount+"</td></tr>");
	}	
    });
    // Create a socket
    socket = new WebSocket('ws://13.78.112.191:20335');
    //socket = new WebSocket('ws://192.168.0.166:20335');
    // Message received on the socket
               
    socket.onopen = function()
    {
      // Web Socket 已连接上，使用 send() 方法发送数据
      GetEventsByHash();
      //alert("Web Socket 已连接上，使用 send() 方法发送数据");
    };
                
    socket.onmessage = function (evt) 
    {
	var received_msg = JSON.parse(evt.data);
    
	if (received_msg["Action"] == "getstorage"){
		result = str2ab(received_msg["Result"])
      		//alert(result);
	}
      GetOrders(1);
    };
                
    socket.onclose = function()
    { 
      // 关闭 websocket
      alert("连接已关闭..."); 
    };
});
