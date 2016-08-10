package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var rootTemplate = template.Must(
	template.New("root").Parse(
		`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<title>Websoket tchat</title>
<style>
.sii-chat{width:400px;padding:10px;background:#ccc;margin:20px auto;}
.sii-chat-content{min-height:400px;max-height:400px;overflow:hidden;overflow-y:scroll;background:#fff;box-shadow:0px 0px 5px 0px #000;margin-bottom:10px;}
.console{min-height:50px;max-height:100px;overflow:hidden;overflow-y:scroll;}
.sii-chat-form input[name="sii-chat-message"]{width:300px;box-shadow:inset 0px 0px 5px 0px #000;}
.sii-chat-form button{width:80px;float:right;color:#ffffff;-moz-box-shadow: 0px 0px 5px #343434;-webkit-box-shadow: 0px 0px 5px #343434;-o-box-shadow: 0px 0px 5px #343434;box-shadow: 0px 0px 5px #343434;-moz-border-radius: 5px;-webkit-border-radius: 5px;border-radius: 5px;border: 1px solid #656565;filter: progid:DXImageTransform.Microsoft.gradient(startColorstr="34cdf9", endColorstr="3166ff"); /* Pour IE seulement et mode gradient à linear */background: -webkit-gradient(linear, left top, left bottom, from(#34cdf9), to(#3166ff));background: -moz-linear-gradient(top center, #34cdf9, #3166ff);}
.sii-chat-form button:active{filter: progid:DXImageTransform.Microsoft.gradient(startColorstr="3166ff", endColorstr="34cdf9"); /* Pour IE seulement et mode gradient à linear */background: -webkit-gradient(linear, left top, left bottom, from(#3166ff), to(#34cdf9));background: -moz-linear-gradient(top center, #3166ff, #34cdf9);}
</style>
</head>
<body>
    <div class="sii-chat">
        <div>Pseudo : <input type="text" name="sii-chat-name" /><button class="sii-chat-login">Valider</button></div>
        <div class="sii-chat-content"></div>
        <div>
            <form class="sii-chat-form" onsubmit="return false;">
                <input type="text" value="" name="sii-chat-message" disabled="disabled"/>
                <button class="sii-chat-send" disabled="disabled">ok</button>
            </form>
        </div>
        <div class="console"></div>
    </div>
<script type="text/javascript">
    var uId = ''; /* pseudo de l'utilisateur*/
    var button = document.getElementsByClassName('sii-chat-send')[0]; /* bouton d'envoi du message */
    var messageInput = document.getElementsByName('sii-chat-message')[0]; /* message à envoyer vers le serveur */
    var buttonUser = document.getElementsByClassName('sii-chat-login')[0]; /* bouton de soumission du pseudo */
    var contentMessage = document.getElementsByClassName('sii-chat-content')[0]; /* div contenant les messages reçus par le serveur*/
    var WebsocketClass = function(host){
        this.host = host
        this.console = document.getElementsByClassName('console')[0];
    };
    WebsocketClass.prototype = {
        initWebsocket : function(){
            var $this = this;
            this.socket = new WebSocket(this.host);
            this.socket.onopen = function(){
                $this.onOpenEvent(this);
            };
            this.socket.onmessage = function(e){
                $this._onMessageEvent(e);
            };
            this.socket.onclose = function(){
                $this._onCloseEvent();
            };
            this.socket.onerror = function(error){
                $this._onErrorEvent(error);
            };
            this.console.innerHTML = this.console.innerHTML + 'websocket init <br />';
        },
        _onErrorEvent :function(err){
            console.log(err);
            this.console.innerHTML = this.console.innerHTML + 'websocket error <br />';
        },
        onOpenEvent : function(socket){
            console.log('socket opened');
            this.console.innerHTML = this.console.innerHTML + 'socket opened Welcome - status ' + socket.readyState + '<br />';
        },
        _onMessageEvent : function(e){
            console.log(e.data)
            msg = JSON.parse(e.data);
            contentMessage.innerHTML = contentMessage.innerHTML + '><strong>' + msg.from + '</strong> : ' + msg.message + '<br />';
            contentMessage.scrollTop = contentMessage.scrollHeight; /* to auto scroll at the bottom */
            this.console.innerHTML = this.console.innerHTML + 'message event lanched <br />';
        },
        _onCloseEvent : function(){
            console.log('connection closed');
            this.console.innerHTML = this.console.innerHTML + 'websocket closed - server not running<br />';
            uId = '';
            document.getElementsByName('sii-chat-name')[0].value = '';
            messageInput.disabled = 'disabled';
            button.disabled = 'disabled';
        },
        sendMessage : function(){
            this.socket.send('{"from":' + JSON.stringify(uId) + ', "message":' + JSON.stringify(messageInput.value) + '}');
            contentMessage.innerHTML = contentMessage.innerHTML + '><strong>' + uId + '</strong> : ' + messageInput.value + '<br />';
            contentMessage.scrollTop = contentMessage.scrollHeight; /* to auto scroll at the bottom */
            messageInput.value = '';
            this.console.innerHTML = this.console.innerHTML + 'websocket message send <br />';
        }
    };
    var socket = new WebsocketClass('ws://{{.}}/socket');
    if(button.addEventListener){
        buttonUser.addEventListener('click', function(e){ /* listen click event on login buton */
            e.preventDefault(); /* stop propagation */
            socket.initWebsocket(); /* init connexion */
            uId = document.getElementsByName('sii-chat-name')[0].value; /* get pseudo */
            messageInput.disabled = ''; /* chat granted  */
            button.disabled = '';
            return false; /* avoid reload on button click */
        }, true);
        button.addEventListener('click',function(e){ /* listen click event on send message button  */
            e.preventDefault();
            socket.sendMessage(); /* send message */
            return false;
        }, true);
    } else{
        console.log('not handled browser');
    }
</script>
</body>
</html>
`,
	),
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	err = rootTemplate.Execute(w, listenAddr)
	if err != nil {
		log.Println(fmt.Sprintf("unable to execute root template: %s", err))
	}
}
