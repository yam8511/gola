package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"gola/internal/bootstrap"
	"sync"
)

/*
| =======================================================
| 以下是加密、解密的部分
| =======================================================
*/

var chip = &struct {
	sync.Once
	block    cipher.Block
	commonIV []byte
}{}

func setupCipher() {
	salt := bootstrap.GetAppConf().App.Salt
	dst := make([]byte, 64)
	base64.StdEncoding.Encode(dst, []byte(salt))

	// 創建加密算法aes
	c, err := aes.NewCipher(dst[:32])
	if err != nil {
		panic(err)
	}
	chip.block = c
	chip.commonIV = dst[32:48]
}

// EncryptSession session加密
func EncryptSession(data []byte) (string, error) {
	chip.Do(setupCipher)

	//加密字符串
	cfb := cipher.NewCFBEncrypter(chip.block, chip.commonIV)
	session := make([]byte, len(data))
	cfb.XORKeyStream(session, data)
	return hex.EncodeToString(session), nil
}

// DecryptSession session解密
func DecryptSession(session string) ([]byte, error) {
	s, err := hex.DecodeString(session)
	if err != nil {
		return nil, err
	}

	chip.Do(setupCipher)

	// 解密字符串
	cfbdec := cipher.NewCFBDecrypter(chip.block, chip.commonIV)
	data := make([]byte, len(s))
	cfbdec.XORKeyStream(data, s)
	return data, nil
}
