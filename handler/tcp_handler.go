package handler

import (
	"io"
	"net"
	//"bytes"
	//"compress/zlib"
	"encoding/binary"
	"github.com/gwtony/gapi/log"
)

// AddHandler urouter udp handler
type URouterHandler struct {
	rh     *RedisHandler
	token  string
	log    log.Log
}

func (handler *URouterHandler) ReadAndParse(conn net.Conn) ([]byte, error) {
	var n, dlen, left uint16
	var headFlag = 1
	head := make([]byte, 3)
	recv := make([]byte, 8192)

	n = 3
	for  {
		if headFlag == 1 {
			s, err := conn.Read(head[3-n:])
			if err != nil {
				if err != io.EOF {
					handler.log.Error("Read failed: %s", err)
				}
				return nil, err
			}
			if uint16(s) < n {
				n -= uint16(s)
				continue
			}
			magic := head[0]
			//handler.log.Debug("magic is %x", magic)
			if magic != 0x77 {
				handler.log.Error("Data magic is invalid")
				return nil, InvalidMagicError
			}
			dlen = binary.LittleEndian.Uint16(head[1:3])
			headFlag = 0
			left = dlen
			continue
		}

		s, err := conn.Read(recv[dlen - left: dlen])
		if err != nil {
			handler.log.Error("Read body failed: %s", err)
			return nil, err
		}

		if uint16(s) < left {
			left -= uint16(s)
			continue
		}

		//handler.log.Debug("parse return %d", dlen)
		return recv[0:dlen], nil
	}
}
func (handler *URouterHandler) ServTcp(conn net.Conn) {
	//var hdata bytes.Buffer
	for {
		//TODO: performance
		frame, err := handler.ReadAndParse(conn)
		if err != nil {
			handler.log.Error("Read and parse failed: %s", err)
			conn.Close()
			break
		}

		uid := string(frame[0:32])
		puid := string(frame[32:64])
		pip := net.IPv4(frame[64], frame[65], frame[66], frame[67]).String()
		pport := binary.BigEndian.Uint16(frame[68:70])
		lip := net.IPv4(frame[70], frame[71], frame[72], frame[73]).String()
		lport := binary.BigEndian.Uint16(frame[74:76])
		hlen := binary.BigEndian.Uint16(frame[76:78])

		//NO ZIP
		//var hdata bytes.Buffer
		//hdata.Reset()
		//w, err := zlib.NewWriterLevel(&hdata, zlib.BestSpeed)
		//w := zlib.NewWriter(&hdata)
		//if err != nil {
		//	handler.log.Error("New zlib writer failed")
		//	//continue
		//	return
		//}
		//w.Write(frame[78:])
		//w.Close()

		//Need not to goroutine
		mpdata, err := EncodeMsgpack(uid, puid, pip, lip, frame[78:], pport, lport, hlen)

		if err != nil {
			handler.log.Error("Encode msgpack failed")
			continue
		}

		handler.rh.RedisSet(uid, mpdata)
	}
}
