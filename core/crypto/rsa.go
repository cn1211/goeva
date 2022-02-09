package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"hash"
	"log"

	"github.com/pkg/errors"
)

type RSACryptMgr struct {
	publicKey  string
	privateKey string
	signAlg    x509.SignatureAlgorithm
}

func NewRSACryptMgr(publicKey, privateKey string, signAlg x509.SignatureAlgorithm) *RSACryptMgr {
	return &RSACryptMgr{
		publicKey:  publicKey,
		privateKey: privateKey,
		signAlg:    signAlg,
	}
}

// base64解码私钥
func (mgr *RSACryptMgr) DecodePrivateKey() ([]byte, error) {
	bPrvKeys, err := base64.StdEncoding.DecodeString(mgr.privateKey)
	if err != nil {
		return nil, errors.New("无法还原私钥")
	}
	return bPrvKeys, nil
}

func (mgr *RSACryptMgr) Sign(data string) (string, error) {
	bPrvKeys, err := mgr.DecodePrivateKey()
	if err != nil {
		return "", err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(bPrvKeys)
	if err != nil {
		log.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}

	if mgr.signAlg != x509.SHA1WithRSA && mgr.signAlg != x509.SHA256WithRSA {
		return "", errors.New("未知的RSA哈希算法")
	}

	var (
		h         hash.Hash
		signature []byte
	)

	if mgr.signAlg == x509.SHA1WithRSA {
		h = crypto.SHA1.New()
		h.Write([]byte(data))
		hashVal := h.Sum(nil)
		signature, err = rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hashVal)
	} else if mgr.signAlg == x509.SHA256WithRSA {
		h = crypto.SHA256.New()
		h.Write([]byte(data))
		hashVal := h.Sum(nil)
		signature, err = rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hashVal)
	}

	if err != nil {
		log.Printf("Error from signing: %s\n", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (mgr *RSACryptMgr) VerifySign(sign, paramsData string) error {
	bsign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	public, _ := base64.StdEncoding.DecodeString(mgr.publicKey)
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}

	var hashVal hash.Hash
	if mgr.signAlg == x509.SHA1WithRSA {
		hashVal = crypto.SHA1.New()
		hashVal.Write([]byte(paramsData))
		err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hashVal.Sum(nil), bsign)
	} else if mgr.signAlg == x509.SHA256WithRSA {
		hashVal = crypto.SHA256.New()
		hashVal.Write([]byte(paramsData))
		err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hashVal.Sum(nil), bsign)
	}

	if err != nil {
		return err
	}
	return nil
}
