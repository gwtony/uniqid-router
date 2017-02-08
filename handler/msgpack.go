package handler
import (
	//"fmt"
	//"bytes"
	"github.com/ugorji/go/codec"
)

func EncodeMsgpack(uid, puid, pip, lip string, data []byte, pport, lport, dlen uint16) ([]byte, error) {
	var b []byte
	var mh codec.MsgpackHandle

	rd := &RouterData{
		Uid: uid,
		Puid: puid,
		Pip: pip,
		Pport: pport,
		Lip: lip,
		Lport: lport,
		Dlen: dlen,
		Data: data,
	}

	enc := codec.NewEncoderBytes(&b, &mh)

	err := enc.Encode(rd)
	if err != nil {
		return nil, err
	}

	//For debug
	//var rh RouterData
	//r := bytes.NewReader(b)
	//dec := codec.NewDecoder(r, &mh)
	//dec.Decode(&rh)
	//fmt.Println("Decode data is", rh)


	return b, nil
}

