package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"regexp"
	"time"
)

func Int2Byte(data int32) (ret []byte) {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(data))
	return buf
}

/*
第三方客户端通过 TCP 协议连接到弹幕服务器(依据指定的 IP 和端口);
第三方接入弹幕服务器列表:
IP 地址:openbarrage.douyutv.com 端口:8601
*/
func connect() (conn *net.TCPConn, err error) {
	fmt.Println("-----*-----DouYu_Spider-----*-----")
	addr, err := net.ResolveTCPAddr("tcp4", "openbarrage.douyutv.com:8601")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn, err = net.DialTCP("tcp4", nil, addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}

func send_msg(conn *net.TCPConn, msg string) {
	var dataLength int32 = int32(len(msg) + 8)
	var code int32 = 689
	byte1 := Int2Byte(dataLength)
	byte2 := Int2Byte(code)

	msgHead := bytes.Join([][]byte{byte1, byte1, byte2}, []byte(""))
	if _, err := conn.Write(msgHead); err != nil {
		fmt.Println(err)
	}
	if _, err := conn.Write([]byte(msg)); err != nil {
		fmt.Println(err)
	}
}

//判断是否为弹幕
func judgeChatmsg(content string) bool {
	reg := regexp.MustCompile("type@=(.*)/rid@")
	matchs := reg.FindStringSubmatch(content)
	if matchs != nil && len(matchs) > 1 && matchs[1] == "chatmsg" {
		return true
	}
	return false
}

func nickNameAndChatMsg(content string) (nickName, chatMsg string) {
	reg := regexp.MustCompile("nn@=([^n]*)/txt@")
	matchs := reg.FindStringSubmatch(content)
	if matchs != nil && len(matchs) > 1 {
		nickName = matchs[1]
	}
	reg = regexp.MustCompile("txt@=([^n]*)/cid@")
	matchs = reg.FindStringSubmatch(content)
	if matchs != nil && len(matchs) > 1 {
		chatMsg = matchs[1]
	}
	return
}

/*
1.客户端向弹幕服务器发送登录请求
2.客户端收到登录成功消息后发送进入弹幕分组请求给弹幕服务器
*/
func danmu(conn *net.TCPConn, room_id string) {
	login := fmt.Sprintf("type@=loginreq/roomid@=%s/\000", room_id)
	send_msg(conn, login)
	joingroup := fmt.Sprintf("type@=joingroup/rid@=%s/gid@=-9999/\000", room_id)
	send_msg(conn, joingroup)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			content := string(buf[:n])
			if judgeChatmsg(content) {
				nick, msg := nickNameAndChatMsg(content)
				fmt.Printf("%s : %s\n", nick, msg)
			}
		}
	}
}

func keepAlive(conn *net.TCPConn) {
	for {
		msg := fmt.Sprintf("type@=keeplive/tick@=%d/\000", time.Now().Unix())
		send_msg(conn, msg)
		time.Sleep(45 * time.Second)
	}
}

func main() {
	conn, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go keepAlive(conn)
	var roomID string
	flag.StringVar(&roomID, "c", "156277", "room_id")
	danmu(conn, roomID)
}
