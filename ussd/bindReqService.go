package ussd

import (
	"bitbucket.org/nndi/dialoguss/configInit"
	bytes2 "bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type BindReq struct {
	Config  *configInit.UssdOpt
	Request *Request
	Conn    net.Conn
	Listen  net.Listener
	wg      *sync.WaitGroup
}

type Request struct {
	HEARDER USSDHEARDER
	BODY    BindReqBody
}

type USSDHEARDER struct {
	Command_Length uint32
	Command_ID     uint32
	Command_Status uint32
	Sender_ID      uint32
	Receiver_ID    uint32
}

type BindReqBody struct {
	System_ID         [11]byte
	Password          [9]byte
	System_Type       [13]byte
	Interface_Version uint32
}

type BindRsp struct {
	HearDer USSDHEARDER
	Body    BindRspBody
}

type BindRspBody struct {
	SystemId [11]byte
}

type SharkReq struct {
	HEARD USSDHEARDER
}

type SharkRsp struct {
	HEADER USSDHEARDER
}

type BeginReq struct {
	HEADER USSDHEARDER
	BODY   beginReqBody
}

type beginReqBody struct {
	Ussd_Version uint8
	Ussd_Op_Type uint8
	MsIsdn       [21]byte
	Service_Code [21]byte
	Code_Scheme  uint8
	Ussd_Content [182]byte
}

type ContinueReq struct {
	BeginReq
}

type ContinueRsp struct {
	BeginReq
}

func (u *BindReq) StartWork(pro string) (err error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("[tcp]recover err:", err)
		}
	}()
	u.wg = new(sync.WaitGroup)

	//系统中所有请求连接超时时间不大于9
	/*	address := fmt.Sprintf("%s:%d", u.Config.Opt.Address, u.Config.Opt.Ports)
		connectionMax := time.Duration(u.Config.Opt.ConnectionMax)*/
	u.Conn, err = net.DialTimeout(pro, "10.71.83.55:4850", time.Millisecond*9000)
	if err != nil {
		log.Printf("[%s]ResolveTCPAddr err =%+v", pro, err)
		return
	}

	bytes := bytes2.NewBuffer(nil)
	err = binary.Write(bytes, binary.LittleEndian, u.Request)
	if nil != err {
		log.Printf("[%s]binary.Write =%+v", pro, err)
		return
	}

	fmt.Printf("bytes.Bytes() =%+v\n", bytes.Bytes())
	//var len int
	_, err = u.Conn.Write(bytes.Bytes())
	if err != nil {
		log.Printf("[%s]Write =%+v", pro, err)
		return
	}
	log.Printf("[%s]********Dial success .u.Conn Pointer =%p", pro, u.Conn)

	/*	err = u.Begin()
		if err != nil {
			return
		}*/

	//心跳
	go func() {
		u.KeepHeart()
	}()

	u.wg.Add(2)
	//保持读取
	go func() {
		u.GetBuffFroMUSSDC()
	}()

	u.wg.Wait()
	return
}

func (u *BindReq) UnBind(conn net.Conn) error {
	return conn.Close()
}

func (u *BindReq) KeepHeart() (err error) {
	defer u.wg.Done()
	header := USSDHEARDER{
		Command_Length: 20,
		Command_ID:     0x00000083,
		Command_Status: 0,
		Sender_ID:      0xFFFFFFFF,
		Receiver_ID:    0xFFFFFFFF,
	}
	bytes := bytes2.NewBuffer(nil)
	err = binary.Write(bytes, binary.LittleEndian, header)
	if nil != err {
		log.Printf("[Shark]binary.Write err=%+v", err)
		return
	}

	for {
		_, err = u.Conn.Write(bytes.Bytes())
		if err != nil {
			log.Printf("[Shark]conn.Write err =%+v", err)
			return
		}

		/*		msg := make([]byte, 1024)
				_, err = u.Conn.Read(msg)
				if err != nil {
					log.Printf("[Shark]Read =%+v", err)
					return
				}

				rsp := &USSDHEARDER{}
				err = binary.Read(bytes2.NewBuffer(msg[:21]), binary.LittleEndian, rsp)
				if nil != err {
					log.Printf("[Shark]Read err = %+v", err)
					return
				}
				log.Printf("[Shark]Read rsp = %+v", rsp)*/

		time.Sleep(time.Second * 5)
	}
}

func (u *BindReq) Begin() (err error) {
	req := BeginReq{
		HEADER: USSDHEARDER{
			Command_Length: 174,
			Command_ID:     0x0000006F,
			Command_Status: 0,
			Sender_ID:      0x01000002,
			Receiver_ID:    0xFFFFFFFF,
		},
		BODY: beginReqBody{
			Ussd_Version: 0x20,
			Ussd_Op_Type: 0x01,
			MsIsdn:       [21]byte{'8', '6', '1', '3', '9', '0', '0', '2', '7', '1', '0', '2', '4'},
			Service_Code: [21]byte{'*', '6', '2', '8', '7', '4', '#'},
			Code_Scheme:  0x0F,
			//Ussd_Content: [182]byte{},
		},
	}

	menu := "1.Borrow Airtime\n2.Borrow Data"
	copy(req.BODY.Ussd_Content[:108], menu)
	//log.Printf("[Begin] copy req.BODY.Ussd_Content  =%+v", req.BODY.Ussd_Content)

	bytes := bytes2.NewBuffer(nil)
	err = binary.Write(bytes, binary.LittleEndian, req)
	if nil != err {
		log.Printf("[Begin]binary.Write err=%+v", err)
		return
	}

	_, err = u.Conn.Write(bytes.Bytes())
	if err != nil {
		log.Printf("[Begin]conn.Write err =%+v", err)
		return
	}
	//log.Printf("[Begin]Success .")

	/*	msg := make([]byte, 1024)
		for {
			req2 := &ContinueReq{}
			_, err = u.Conn.Read(msg)
			if err != nil {
				log.Printf("[Begin]Read =%+v", err)
				return
			}

			//second menu
			err = binary.Read(bytes2.NewBuffer(msg), binary.LittleEndian, req2)
			if nil != err {
				log.Printf("[Begin]Read err = %+v", err)
				return
			}
			log.Printf("[Begin]Read rsp = %+v", req)

			switch req.BODY.Ussd_Op_Type {
			case 1:
				//二级菜单
				log.Printf("[Begin] 1:*************** ussdOpType = %d", 1)
				req2.HEADER.Command_Length = 186
				req2.HEADER.Command_ID = 0x00000070
				req2.BODY.Ussd_Version = 0x20
				req2.BODY.Ussd_Op_Type = 0x01
				req2.BODY.Code_Scheme = 0x0F
				menu := "1.choose loan amount\n2.check eligibility\n3.check debt status\n4. enable loans\n5.Disable Loans"
				copy(req.BODY.Ussd_Content[:182], menu)
			case 2:
				//三级菜单
				log.Printf("[Begin] 3:*************** ussdOpType = %+v", 3)
				req.HEADER.Command_Length = 186
				req.HEADER.Command_ID = 0x00000070
				req.BODY.Ussd_Version = 0x20
				req.BODY.Ussd_Op_Type = 0x01
				menu := "1.Get N50\n2.Get N100\n3.Get N200\n4.Get N500\n5.Get N1000\n6.Get N2000"
				copy(req.BODY.Ussd_Content[:182], menu)

			case 3:
				//三级菜单
				log.Printf("[Begin] 3:*************** ussdOpType = %+v", 3)
				req.HEADER.Command_Length = 186
				req.HEADER.Command_ID = 0x00000070
				req.BODY.Ussd_Version = 0x20
				req.BODY.Ussd_Op_Type = 0x01
				menu := "1.Get N50\n2.Get N100\n3.Get N200\n4.Get N500\n5.Get N1000\n6.Get N2000"
				copy(req.BODY.Ussd_Content[:182], menu)

			case 4:
				//三级菜单
				log.Printf("[Begin] 3:*************** ussdOpType = %+v", 3)
				req.HEADER.Command_Length = 186
				req.HEADER.Command_ID = 0x00000070
				req.BODY.Ussd_Version = 0x20
				req.BODY.Ussd_Op_Type = 0x01
				menu := "1.Get N50\n2.Get N100\n3.Get N200\n4.Get N500\n5.Get N1000\n6.Get N2000"
				copy(req.BODY.Ussd_Content[:182], menu)
			default:
				log.Printf("[Begin] default:*************** ussdOpType = %+v", 3)

			}

			bytes := bytes2.NewBuffer(nil)
			err = binary.Write(bytes, binary.LittleEndian, req)
			if nil != err {
				log.Printf("[Begin]binary.Write err=%+v", err)
				return
			}

			_, err = u.Conn.Write(bytes.Bytes())
			if err != nil {
				log.Printf("[Begin]conn.Write err =%+v", err)
				return
			}
		}*/

	return
}

func (u *BindReq) GetBuffFroMUSSDC() (err error) {
	defer u.wg.Done()

	for {
		msg := make([]byte, 1024)
		var rlen int
		if err := u.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3)); err != nil {
			return err
		}
		rlen, err = u.Conn.Read(msg)
		if err != nil {
			log.Printf("[GetBuffFroMUSSDC] Read err =%+v", err)
			break
		}
		if err := u.Conn.SetReadDeadline(time.Time{}); err != nil {
			return err
		}

		header := &USSDHEARDER{}
		err = binary.Read(bytes2.NewBuffer(msg[:21]), binary.LittleEndian, header)
		switch header.Command_ID {
		case 103:
			//绑定
			rsp := &BindRsp{}
			err = binary.Read(bytes2.NewBuffer(msg), binary.LittleEndian, rsp)
			if nil != err {
				log.Printf("[BIAND]Read err = %+v", err)
				continue
			}
			log.Printf("[BIAND]Read rsp = %+v", rsp)
		case 111:
			//begin
			rsp := &BeginReq{}
			err = binary.Read(bytes2.NewBuffer(msg), binary.LittleEndian, rsp)
			if nil != err {
				log.Printf("[begin]Read err = %+v", err)
				continue
			}
			//u.doMenu(rsp)
			u.doBegin(rsp)
		case 112:
			//二级菜单请求
			rsp := &BeginReq{}
			err = binary.Read(bytes2.NewBuffer(msg), binary.LittleEndian, rsp)
			if nil != err {
				log.Printf("[menu_2]Read err = %+v", err)
				continue
			}
			u.doMenu(rsp)
		case 114:
			//abort
			rsp := &USSDHEARDER{}
			err = binary.Read(bytes2.NewBuffer(msg), binary.LittleEndian, rsp)
			if nil != err {
				log.Printf("[abort]Read err = %+v", err)
				continue
			}
			log.Printf("[abort]Read rsp = %+v", rsp)
			//u.doMenu(rsp)
		case 132:
			//心跳
			rsp := &SharkReq{}
			err = binary.Read(bytes2.NewBuffer(msg), binary.LittleEndian, rsp)
			if nil != err {
				log.Printf("[Shark]Read err = %+v", err)
				continue
			}
			log.Printf("[Shark]Read rsp = %+v", rsp)
		default:
			log.Printf("[GetBuffFroMUSSDC] Con't find u.Conn Pointer = %p, comand =%+v ,readlen =%d", u.Conn, header, rlen)
		}

		//time.Sleep(1 * time.Second)
	}

	return
}

func (u *BindReq) doMenu(req *BeginReq) (err error) {
	log.Printf("[doMenu] req =%+v", req)

	rsp := &BeginReq{}
	ussdContent := string(req.BODY.Ussd_Content[:])
	log.Printf("[doMenu] ussdContent =%s",ussdContent)
	switch req.BODY.Ussd_Content[0] {
	case 49:
		//二级菜单
		log.Printf("[doMenu] 1:*************** ussdContent = %+v ,rsp.BODY.MsIsdn =%+v", ussdContent, rsp)
		rsp.HEADER.Command_Length = 53
		rsp.HEADER.Command_ID = 0x00000070
		//rsp.HEADER.Command_Status = 0
		rsp.HEADER.Sender_ID = req.HEADER.Sender_ID
		rsp.HEADER.Receiver_ID = req.HEADER.Sender_ID

		rsp.BODY.Ussd_Version = 0x20
		rsp.BODY.Ussd_Op_Type = 0x30
		rsp.BODY.MsIsdn = req.BODY.MsIsdn
		rsp.BODY.Code_Scheme = req.BODY.Code_Scheme
		rsp.BODY.Service_Code = req.BODY.Service_Code
		menu := "1.choose loan amount\n2.check eligibility\n3.check debt status\n4. enable loans\n5.Disable Loans"
		copy(rsp.BODY.Ussd_Content[:182], menu)


	case 50:
		//三级菜单
		log.Printf("[doMenu] 3:*************** ussdContent = %+v ,rsp.BODY.MsIsdn =%+v", ussdContent, rsp)
		rsp.HEADER.Command_Length = 186
		rsp.HEADER.Command_ID = 0x00000070
		//rsp.HEADER.Command_Status = 0
		rsp.HEADER.Sender_ID = req.HEADER.Sender_ID
		rsp.HEADER.Receiver_ID = req.HEADER.Sender_ID

		rsp.BODY.Ussd_Version = 0x20
		rsp.BODY.Ussd_Op_Type = 0x01
		rsp.BODY.MsIsdn = req.BODY.MsIsdn
		rsp.BODY.Code_Scheme = req.BODY.Code_Scheme
		rsp.BODY.Service_Code = req.BODY.Service_Code
		menu := "1.Get N50\n2.Get N100\n3.Get N200\n4.Get N500\n5.Get N1000\n6.Get N2000"
		copy(rsp.BODY.Ussd_Content[:182], menu)

	case 51:
		//三级菜单
		menu := "1.Get N50\n2.Get N100\n3.Get N200\n4.Get N500\n5.Get N1000\n6.Get N2000"
		copy(rsp.BODY.Ussd_Content[:182], menu)
	default:
		u.doBegin(rsp)
	}

	rsp.HEADER.Command_Length = 186
	rsp.HEADER.Command_ID = 0x00000070
	rsp.BODY.Ussd_Version = 0x20
	rsp.BODY.Ussd_Op_Type = 0x01

	log.Printf("[doMenu] rsp =%+v", rsp)
	bytes := bytes2.NewBuffer(nil)
	err = binary.Write(bytes, binary.LittleEndian, rsp)
	if nil != err {
		log.Printf("[doMenu]binary.Write err=%+v", err)
		return
	}

	_, err = u.Conn.Write(bytes.Bytes())
	if err != nil {
		log.Printf("[doMenu]conn.Write err =%+v", err)
		return
	}
	return
}

func (u *BindReq) doBegin(req *BeginReq) (err error) {
	rsp := &BeginReq{}

	rsp.HEADER.Command_Length = 186
	rsp.HEADER.Command_ID = 0x00000070
	rsp.BODY.Ussd_Version = 0x20
	rsp.BODY.Ussd_Op_Type = 0x01
	rsp.HEADER.Sender_ID = req.HEADER.Sender_ID
	rsp.HEADER.Receiver_ID = req.HEADER.Sender_ID
	rsp.BODY.MsIsdn = req.BODY.MsIsdn
	rsp.BODY.Service_Code = req.BODY.Service_Code
	menu := "1.Borrow Airtime\n2.Borrow Data"
	copy(rsp.BODY.Ussd_Content[:182], menu)

	bytes := bytes2.NewBuffer(nil)
	err = binary.Write(bytes, binary.LittleEndian, rsp)
	if nil != err {
		log.Printf("[doMenu]binary.Write err=%+v", err)
		return
	}

	_, err = u.Conn.Write(bytes.Bytes())
	if err != nil {
		log.Printf("[doMenu]conn.Write err =%+v", err)
		return
	}
	return
}
