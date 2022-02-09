package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/pkg/errors"
)

type AESCryptMgr struct {
	key []byte
}

func NewAESCryptMgr(key string) *AESCryptMgr {
	return &AESCryptMgr{
		key: padding([]byte(key)),
	}
}

func (mgr *AESCryptMgr) Encrypt(srcData []byte) ([]byte, error) {
	// AES CBC对被加密数据长度有要求, key同
	srcData = padding(srcData)

	block, err := aes.NewCipher(mgr.key)
	if err != nil {
		return nil, err
	}

	encryptData := make([]byte, len(srcData))

	// key直接作为iv
	mode := cipher.NewCBCEncrypter(block, mgr.key)
	mode.CryptBlocks(encryptData, srcData)
	return encryptData, nil
}

func (mgr *AESCryptMgr) Decrypt(srcData []byte) ([]byte, error) {
	// 密文字节数组长度必须是aes.BlockSize的整数倍
	if len(srcData)%aes.BlockSize != 0 {
		return nil, errors.New("cipher text is not correct!")
	}

	block, err := aes.NewCipher(mgr.key)
	if err != nil {
		return nil, err
	}

	// key直接作为iv
	mode := cipher.NewCBCDecrypter(block, mgr.key)

	dstContentData := make([]byte, len(srcData))
	mode.CryptBlocks(dstContentData, srcData)

	// 解密后,移除字节数组末尾的额外\0字节(加密时padding填充的)
	dstContentData = bytes.TrimRight(dstContentData, "\x00")

	return dstContentData, nil
}

// 加密字符串为b64结果返回
func (mgr *AESCryptMgr) EncryptToBase64(src string) (string, error) {
	srcData := []byte(src)
	dst, err := mgr.Encrypt(srcData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dst), nil
}

// 从b64加密的加密字符串解码返回
func (mgr *AESCryptMgr) DecryptFromBase64(src string) (string, error) {
	srcData, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	dst, err := mgr.Decrypt(srcData)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}

// AES/CBC/NoPadding
// https://github.com/golang/go/issues/24402#issuecomment-373299954
func padding(data []byte) []byte {
	if len(data)%aes.BlockSize == 0 {
		return data
	}

	needPadLen := aes.BlockSize - len(data)%aes.BlockSize
	for i := 0; i < needPadLen; i++ {
		data = append(data, 0)
	}
	return data
}
