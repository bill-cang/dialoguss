package main

import (
	"bitbucket.org/nndi/dialoguss/ussd"
	"fmt"
	"log"
)

func init() {

}

func main() {
/*	ussdConfig, err := mod.GetUSSDConfig()
	if err != nil {
		return
	}*/

	pros := []string{"tcp"}
	for _, v := range pros {
		log.Printf("[%s] begin ===============================================>>>>>>", v)
		var u = ussd.BindReq{
			Config: nil,
			Request: &ussd.Request{
				HEARDER: ussd.USSDHEARDER{
					Command_Length: 57,
					Command_ID:     0x00000065,
					Command_Status: 0,
					Sender_ID:      0xFFFFFFFF,
					Receiver_ID:    0xFFFFFFFF,
				},
				BODY: ussd.BindReqBody{
					System_ID: [11]byte{'a', 'p', 'p', 'l', 'i', 'n', 'k'},
					Password:  [9]byte{'a', 'p', 'p', 'l', 'i', 'n', 'k'},
					//System_Type: [13]byte{},

					/*				System_ID:         [11]byte{'t', 'e', 's', 't'},
									Password:          [9]byte{'t', 'e', 's', 't'},
									System_Type:       [13]byte{'U', 'S', 'S', 'D'},*/
					Interface_Version: 0x00000010,
				},
			},
		}

		err := u.StartWork(v)
		if err != nil {
			fmt.Printf("StartWork err =%+v", err)
			return
		}
	}
}
