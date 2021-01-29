package cmdb

import (
	"encoding/json"
	"go-xops/dto/service/terminal/sftp"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type sftp_resp struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func Sftp_ssh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer wsConn.Close()

	var auth WsAuth

	if c.ShouldBindUri(&auth) != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("参数错误\r\n"))
		wsConn.Close()
		return
	}

	for {
		_, _, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			wsConn.Close()
			return
		}

		resp_msg := sftp_resp{}

		path, err := sftp.Client.C[auth.Sid].Sftp.Getwd()
		if err != nil {
			resp_msg.Code = 400
			resp_msg.Type = "connect"
			resp_msg.Msg = "SFTP连接失败"
			msg, _ := json.Marshal(resp_msg)
			if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("sftp connect err:", err)
			}
			return
		}

		resp_msg.Code = 200
		resp_msg.Type = "connect"
		resp_msg.Msg = "连接成功"
		resp_msg.Data = path
		msg, _ := json.Marshal(resp_msg)
		if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("sftp return err:", err)
			return
		}

		break
		//break
	}
	quitChan := make(chan bool, 2)
	go sftp.Client.C[auth.Sid].ReceiveWsMsg(wsConn, quitChan)
	<-quitChan

}
