package ussd

import (
	bytes2 "bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"
)

/*type SharkReq struct {
	HEARD USSDHEARDER
}

type SharkRsp struct {
	HEADER USSDHEARDER
}

type SharkRspHeader struct {
	Command_Length uint32
	Command_ID     uint32
}

type SharkRspBody struct {
	Command_Status uint32
	Sender_ID      uint32
	Receiver_ID    uint32
}
*/
func KeepHeart(conn net.Conn, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	header := USSDHEARDER{
		Command_Length: 20,
		Command_ID:     0x00000083,
		Command_Status: 0,
		Sender_ID:      0xFFFFFFFF,
		Receiver_ID:    0xFFFFFFFF,
	}
	bytes := bytes2.NewBuffer(nil)

	for {
		err = binary.Write(bytes, binary.LittleEndian, header)
		if nil != err {
			log.Printf("[Shark]binary.Write err=%+v", err)
			return
		}

		_, err = conn.Write(bytes.Bytes())
		if err != nil {
			log.Printf("[Shark]conn.Write err =%+v", err)
			return
		}

		time.Sleep(time.Second * 3)

		msg := make([]byte, 1024)
		_, err = conn.Read(msg)
		if err != nil {
			log.Printf("[Shark]Read =%+v", err)
			return
		}
		rsp := &SharkRsp{}
		err = binary.Read(bytes2.NewBuffer(msg), binary.LittleEndian, rsp)
		if nil != err {
			log.Printf("[Shark]Read err = %+v", err)
			return
		}
		log.Printf("[Shark]Read rsp = %+v", rsp)

		time.Sleep(time.Second * 2)
	}
}
