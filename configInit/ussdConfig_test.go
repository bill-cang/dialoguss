package configInit

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestInitUSSDConfig(t *testing.T) {
	config, err := GetUSSDConfig()
	if err != nil {
		fmt.Printf("err = %+v", err)
		return
	}
	fmt.Printf("config = %+v\n", config)

	duration := time.Duration(9000) * time.Millisecond
	fmt.Printf("duration =%+v\n", duration)

	/*	bytes := []byte{0x39, 00, 00, 00, 0x65, 00, 00, 00, 00, 00, 00, 00, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0x74, 0x65, 0x73, 0x74, 00, 00, 00, 00, 00, 00, 00, 0x74,
				0x65, 0x73, 0x74, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00,
				00, 00, 00, 00, 00, 0x10, 00, 00, 00}

			bytes = []byte{0x1f, 00, 00, 00, 0x67, 00, 00, 00, 0x01, 0x08, 00, 00, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0x74, 0x65, 0x73, 0x74, 00, 00, 00, 00, 00, 00, 00}

		bytes := []byte{
			0x39, 00, 00, 00, 0x65, 00, 00, 00,
			00, 00, 00, 00, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x61, 0x70, 0x70, 0x6c,
			0x69, 0x6e, 0x6b, 00, 00, 00, 00, 0x61,
			0x70, 0x70, 0x6c, 0x69, 0x6e, 0x6b, 00, 00,
			0x55, 0x53, 0x53, 0x44, 00, 00, 00, 00,
			00, 00, 00, 00, 00, 0x10, 00, 00,
			00}

		for _, v := range bytes {
			str := string(v)
			fmt.Printf("%s ", str)
		}*/

	/*	hex := "0x0000006F"
		val := hex[2:]

		n, err := strconv.ParseUint(val, 16, 32)
		if err != nil {
			panic(err)
		}
		//n2 := uint32(n)
		fmt.Println(n)

		menu := "1.Borrow Airtime\n2.Borrow Data ckx"
		var Ussd_Content [182]byte
		copy(Ussd_Content[:182], menu)
		fmt.Printf("Ussd_Content = %+v\n", Ussd_Content)

		fmt.Println("*************************************************************************")
		msi := [182]byte{}
		msi[0] = 49
		str := string(msi[:])
		fmt.Printf("str = %+v , len = %d\n", str, len(str))
		if str[:1] == "\x00" {
		}*/

	msi2 := []byte{49, 0, 0, 0, 0, 0, 0}
	//Ussd_Content:[49 46 99 104 111 111 115 101 32 108 111 97 110 32 97 109 111 117 110 116 10 50 46 99 104 101 99 107 32 101 108 105 103 105 98 105 108 105 116 121 10 51 46 99 104 101 99 107 32 100 101 98 116 32 115 116 97 116 117 115 10 52 46 32 101 110 97 98 108 101 32 108 111 97 110 115 10 53 46 68 105 115 97 98 108 101 32 76 111 97 110 115 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	str2 := string(msi2)
	fmt.Printf("str2 = %s , len2 = %d ,msi2[0] =%d\n", str2, len(str2), msi2[0])
	str2 = strings.Replace(str2, "\x00", "", -1)
	atoi, err := strconv.Atoi(str2)
	if err != nil {
		fmt.Printf("err =%+v", err)
		return
	}
	fmt.Printf("atoi =%d\n", atoi)

}
