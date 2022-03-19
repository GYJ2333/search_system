package storage

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/GYJ2333/search_system/log"
	"github.com/valyala/gozstd"
)

type SProxy struct {
	rootPath string
}

func (sp *SProxy) Init(storagePath string) {
	sp.rootPath = storagePath
}

func (sp *SProxy) Write(key string, value []byte) (err error) {
	f, err := os.OpenFile(sp.rootPath+key, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0766)
	defer func() {
		err = f.Close()
		if err != nil {
			log.StorageLogger.Printf("Close file(%s) err(%v)", key, err)
		}
	}()
	if err != nil {
		log.StorageLogger.Printf("Open file(%s) err(%v)", key, err)
		return fmt.Errorf("open file(%s) err(%v)", key, err)
	}
	_, err = f.Write(value)
	if err != nil {
		log.StorageLogger.Printf("Write file(%s) err(%v)", key, err)
		return fmt.Errorf("write file(%s) err(%v)", key, err)
	}
	return nil
}

func (sp *SProxy) Read(key string) ([]byte, error) {
	rowData, err := ioutil.ReadFile(sp.rootPath + key)
	if err != nil {
		log.StorageLogger.Printf("Read file(%s) err(%v)", key, err)
		return nil, fmt.Errorf("read file(%s) err(%v)", key, err)
	}

	decompressedData := make([]byte, len(rowData)*2)
	decompressedData, err = gozstd.Decompress(decompressedData, rowData)
	if err != nil {
		log.StorageLogger.Printf("Decompress file(%s) err(%v)", key, err)
		return nil, fmt.Errorf("decompress file(%s) err(%v)", key, err)
	}
	return decompressedData, nil
}

func (sp *SProxy) Delete() error {
	return nil
}
