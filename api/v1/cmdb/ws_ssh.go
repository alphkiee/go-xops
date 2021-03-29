package cmdb

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-xops/internal/service/cmd"
	"go-xops/internal/service/terminal"
	"go-xops/internal/service/terminal/sftp"
	"go-xops/pkg/utils"

	"go-xops/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type SerInfo struct {
	ID       uint
	Ip       string
	Port     int
	Username string
	Password string
	BindUser uint
}

type WsAuth struct {
	Sid string `uri:"sid" binding:"required,uuid"`
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func jsonError(c *gin.Context, msg interface{}) {
	c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": msg})
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		//logrus.WithError(err).Error("gin context http handler error")
		jsonError(c, err.Error())
		return true
	}
	return false
}

func wshandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		logrus.WithError(err).Error("handler ws ERROR:")
		dt := time.Now().Add(time.Second)
		if err := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			logrus.WithError(err).Error("websocket writes control message failed:")
		}
		return true
	}
	return false
}

func WsSsh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if handleError(c, err) {
		return
	}
	defer wsConn.Close()

	cols, err := strconv.Atoi(c.DefaultQuery("cols", "120"))
	if wshandleError(wsConn, err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "32"))
	if wshandleError(wsConn, err) {
		return
	}

	hostId := utils.Str2Uint(c.Query("host_id"))
	host, err := cmd.GetHostByid(hostId)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	var auth WsAuth
	for {
		_, _, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			wsConn.Close()
			//logrus.WithError(err).Error("reading webSocket message failed")
			return
		}
		cache, _ := cache.New(time.Second * 15)
		err = cache.GetCache(auth.Sid)

		if err != nil {
			wsConn.WriteMessage(websocket.TextMessage, []byte("连接超时，请重试！\r\n"))
			wsConn.Close()
			return
		}
		break

	}

	client, err := terminal.NewSSHClient(&host)
	if wshandleError(wsConn, err) {
		return
	}
	defer client.Close()

	ssConn, err := terminal.NewSshConn(cols, rows, client)

	if wshandleError(wsConn, err) {
		return
	}
	sftp.Client.Lock()
	sftp.Client.C[auth.Sid] = &sftp.MyClient{ssConn.SftpClient}
	sftp.Client.Unlock()
	defer func() {
		sftp.Client.Lock()
		delete(sftp.Client.C, auth.Sid) //释放SFTP客户端
		sftp.Client.Unlock()
	}()
	defer ssConn.Close()
	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan
}
