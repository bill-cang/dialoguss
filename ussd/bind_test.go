package ussd

import (
	"testing"
)

/*
0000   39 00 00 00 65 00 00 00 00 00 00 00 ff ff ff ff   9...e...........
0010   ff ff ff ff 74 65 73 74 00 00 00 00 00 00 00 74   ....test.......t
0020   65 73 74 00 00 00 00 00 00 00 00 00 00 00 00 00   est.............
0030   00 00 00 00 00 10 00 00 00                        .........

0000   39 00 00 00 65 00 00 00 00 00 00 00 ff ff ff ff   9...e...........
0010   ff ff ff ff 74 65 73 74 00 00 00 00 00 00 00 74   ....test.......t
0020   65 73 74 00 00 00 00 00 00 00 00 00 00 00 00 00   est.............
0030   00 00 00 00 00 10 00 00 00                        .........

[57 0 0 0 101 0 0 0 0 0 0 0 255 255 255 255 255 255 255 255 116 101 115 116 0 0 0 0 0 0 0 116 101 115 116 0 0 0 0 0 85 83 83 68 0 0 0 0 0 0 0 0 0 16 0 0 0]



0000   39 00 00 00 65 00 00 00 00 00 00 00 ff ff ff ff   9...e...........
0010   ff ff ff ff 61 70 70 6c 69 6e 6b 00 00 00 00 61   ....applink....a
0020   70 70 6c 69 6e 6b 00 00 55 53 53 44 00 00 00 00   pplink..USSD....
0030   00 00 00 00 00 10 00 00 00                        .........



0000   00 00 00 39 00 00 00 65 00 00 00 00 ff ff ff ff   ...9...e........
0010   ff ff ff ff 61 70 70 6c 69 6e 6b 00 00 00 00 61   ....applink....a
0020   70 70 6c 69 6e 6b 00 00 55 53 53 44 00 00 00 00   pplink..USSD....
0030   00 00 00 00 00 00 00 00 10                        .........


0000   00 00 00 39 00 00 00 65 00 00 00 00 ff ff ff ff   ...9...e........
0010   ff ff ff ff 61 70 70 6c 69 6e 6b 00 00 00 00 61   ....applink....a
0020   70 70 6c 69 6e 6b 00 00 55 53 53 44 00 00 00 00   pplink..USSD....
0030   00 00 00 00 00 00 00 00 10                        .........

*/

func TestUSSD_Bind(t *testing.T) {

	var u = BindReq{
		Config: nil,
		Request: &Request{
			HEARDER: USSDHEARDER{
				Command_Length: 57,
				Command_ID:     0x00000065,
				Command_Status: 0,
				Sender_ID:      0xFFFFFFFF,
				Receiver_ID:    0xFFFFFFFF,
			},
			BODY: BindReqBody{
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

	c :=&ContinueRsp{BeginReq{
		BODY: beginReqBody{
			Ussd_Op_Type: 3,
		},
	}}

	u.doMenu(c)

}
