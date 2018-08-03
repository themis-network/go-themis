package escrow

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

var logger *log.Logger

//字节转换成整形
func BytesToUint32(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return tmp
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func init() {
	fileName := "escrow.log"
	logFile,err  := os.Create(fileName)
	if err != nil {
		logger.Fatalln("open file error !")
	}
	// 创建一个日志对象
	logger = log.New(logFile,"",log.LstdFlags)
}