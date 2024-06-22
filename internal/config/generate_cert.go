package config

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/augustjourney/urlshrt/internal/logger"
)

func (c *Config) generateCerts() error {

	cert := &x509.Certificate{

		SerialNumber: big.NewInt(1658),

		Subject: pkix.Name{
			Organization: []string{"URLSHRT"},
			Country:      []string{"RU"},
		},

		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},

		NotBefore: time.Now(),

		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.Log.Fatal(err)
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		logger.Log.Fatal(err)
		return err
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	pemFile, err := os.Create(c.CertPemPath)
	if err != nil {
		logger.Log.Fatal(err)
		return err
	}

	defer pemFile.Close()

	_, err = pemFile.Write(certPEM.Bytes())
	if err != nil {
		logger.Log.Fatal(err)
		return err
	}

	keyFile, err := os.Create(c.CertKeyPath)
	if err != nil {
		logger.Log.Fatal(err)
		return err
	}

	defer keyFile.Close()

	_, err = keyFile.Write(privateKeyPEM.Bytes())
	if err != nil {
		logger.Log.Fatal(err)
		return err
	}

	return nil
}

func (c *Config) checkIfCertsExist() bool {
	_, err1 := os.Stat(c.CertKeyPath)
	_, err2 := os.Stat(c.CertPemPath)

	return err1 == nil && err2 == nil
}

// Проверяет — есть ли уже сгенерированные сертификаты
// Если нет — генерирует новые и возвращает пути на них
// Если уже есть — возвращает пути на них
func (c *Config) GetCerts() (string, string, error) {
	exist := c.checkIfCertsExist()
	if exist {
		logger.Log.Info("Certs exist")
		return c.CertPemPath, c.CertKeyPath, nil
	}

	logger.Log.Info("Certs do not exist, generating new ones")

	err := c.generateCerts()
	if err != nil {
		return "", "", err
	}

	return c.CertPemPath, c.CertKeyPath, nil
}
