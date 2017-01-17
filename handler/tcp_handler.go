package handler

import (
	"io"
	"net"
	"bytes"
	//"errors"
	//"fmt"
	//"time"
	//"math/rand"
	//"strings"
	//"strconv"
	//"io/ioutil"
	//"net/http"
	//"encoding/json"
	"compress/zlib"
	"encoding/binary"
	//"encoding/hex"
	//"bytes"
	//"gopkg.in/redis.v5"
	//"github.com/ugorji/go/codec"
	"github.com/gwtony/gapi/log"
	//"github.com/gwtony/gapi/api"
	//"github.com/gwtony/gapi/errors"
)

// AddHandler urouter udp handler
type URouterHandler struct {
	rh     *RedisHandler
	//domain string
	token  string
	log    log.Log
}

func (handler *URouterHandler) ReadAndParse(conn net.Conn) ([]byte, error) {
	var n, dlen, left uint16
	var headFlag = 1
	head := make([]byte, 3)
	recv := make([]byte, 4096)

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
				handler.log.Error("data magic is invalid")
				return nil, InvalidMagicError
			}
			dlen = binary.LittleEndian.Uint16(head[1:3])
			headFlag = 0
			left = dlen
			continue
		}

		//handler.log.Debug("Dlen is %d, left is %d", dlen, left)
		s, err := conn.Read(recv[dlen - left: dlen])
		if err != nil {
			handler.log.Error("Read body failed: %s", err)
			return nil, err
		}
		//handler.log.Debug("Read %d data", s)
		if uint16(s) < left {
			left -= uint16(s)
			//handler.log.Debug("after left is %d", left)
			continue
		}
		return recv[0:dlen], nil
	}
}
func (handler *URouterHandler) ServTcp(conn net.Conn) {
	var hdata bytes.Buffer
	//pos := 0

	//data := make([]byte, 4096)
	for {
		//TODO: performance
		frame, err := handler.ReadAndParse(conn)
		if err != nil {
			handler.log.Error("Read and parse failed: %s", err)
			conn.Close()
			break
		}
		//handler.log.Debug("frame is %s", frame)

		uid := string(frame[0:32])
		puid := string(frame[32:64])
		pip := net.IPv4(frame[64], frame[65], frame[66], frame[67]).String()
		pport := binary.BigEndian.Uint16(frame[68:70])
		lip := net.IPv4(frame[70], frame[71], frame[72], frame[73]).String()
		lport := binary.BigEndian.Uint16(frame[74:76])
		hlen := binary.BigEndian.Uint16(frame[76:78])

		//handler.log.Debug("data is ", string(data[81:81+hlen]))

		//handler.log.Debug("uid: %s, puid: %s, pip: %s, pport: %d, lip: %s, lport: %d, dlen: %d", uid, puid, pip, pport, lip, lport, hlen)

		hdata.Reset()
		w, err := zlib.NewWriterLevel(&hdata, zlib.BestSpeed)
		if err != nil {
			handler.log.Error("New zlib writer failed")
			continue
		}
		w.Write(frame[78:78+hlen])
		w.Close()

		//handler.log.Debug("hdata is %s", string(hdata.Bytes()))
		mpdata, err := EncodeMsgpack(uid, puid, pip, lip, hdata.Bytes(), pport, lport, hlen)
		if err != nil {
			handler.log.Error("Encode msgpack failed")
			continue
		}

		err = handler.rh.Set(uid, mpdata, UROUTER_DEFAULT_TTL)
		if err != nil {
			handler.log.Error("Set to Redis failed")
		}
	}
}
