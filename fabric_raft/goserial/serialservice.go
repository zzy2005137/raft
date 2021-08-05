//01 04 01 00 00 08 f0 30 
package goserial

import (


        "github.com/tarm/serial"
	"fmt"


)
/*
func main() {
        d,err:=Serial("/dev/ttyUSB0",9600)
        
        if err != nil{
                fmt.Println(err)
        }else{
                fmt.Println(d)
        }

        
     

}
*/

func Serial (com string, baud int, /*time int*/) (float64, error) {

        c := &serial.Config{Name: com, Baud: baud/*, ReadTimeout: time.Second * time*/}
        s, err := serial.OpenPort(c)
        if err != nil {
                fmt.Println("串口打开失败")
                return -1, err
        }
        
        by := []byte { 0x01, 0x04, 0x01, 0x00, 0x00, 0x08, 0xf0, 0x30 }
        _, err = s.Write( by )
        if err != nil {
                fmt.Println("串口写入失败")
                return -1, err
        }
        
        buf := make([]byte, 128)

        //_返回数据长度
        
        _, err = s.Read(buf)
        if err != nil {
                fmt.Println("串口读取失败")
                return -1, err
        }
        

        d1 := float64(buf[3])
        d2 := float64(buf[4])
        //receive:=buf[3:5]
        receive := (d1*256+d2)*20/4095



        return receive, nil
}
