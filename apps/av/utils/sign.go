package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func GenerateCert(domain string) {
	var err error
	rootKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logrus.Info("[-] ", err.Error())
		os.Exit(0)
	}
	certs, err := GetCertificatesPEM(domain + ":443")
	if err != nil {
		logrus.Info("[-] ", err.Error())
		os.Exit(0)
	}
	block, _ := pem.Decode([]byte(certs))
	cert, _ := x509.ParseCertificate(block.Bytes)

	keyToFile(domain+".key", rootKey)

	SubjectTemplate := x509.Certificate{
		SerialNumber: cert.SerialNumber,
		Subject: pkix.Name{
			CommonName: cert.Subject.CommonName,
		},
		NotBefore:             cert.NotBefore,
		NotAfter:              cert.NotAfter,
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	IssuerTemplate := x509.Certificate{
		SerialNumber: cert.SerialNumber,
		Subject: pkix.Name{
			CommonName: cert.Issuer.CommonName,
		},
		NotBefore: cert.NotBefore,
		NotAfter:  cert.NotAfter,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &SubjectTemplate, &IssuerTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		logrus.Info("[-] ", err.Error())
		os.Exit(0)
	}
	certToFile(domain+".pem", derBytes)
}

func keyToFile(filename string, key *rsa.PrivateKey) {
	file, err := os.Create(filename)
	if err != nil {
		logrus.Info("[-] Unable to marshal Create private key: ", err.Error())
		os.Exit(0)
	}
	defer file.Close()
	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		logrus.Info("[-] Unable to marshal RSA private key: ", err.Error())
		os.Exit(0)
	}
	if err := pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}); err != nil {
		logrus.Info("[-] Unable to pem RSA private key: ", err.Error())
		os.Exit(0)
	}
}

func certToFile(filename string, derBytes []byte) {
	certOut, err := os.Create(filename)
	if err != nil {
		logrus.Info("[-] Failed to Open cert.pem for Writing: ", err.Error())
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		logrus.Info("[-] Failed to Write Data to cert.pem: ", err.Error())
	}
	if err := certOut.Close(); err != nil {
		logrus.Info("[-] Error Closing cert.pem: ", err.Error())
	}
}

func GetCertificatesPEM(address string) (string, error) {
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", err
	}
	defer conn.Close()
	var b bytes.Buffer
	for _, cert := range conn.ConnectionState().PeerCertificates {
		err := pem.Encode(&b, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func GeneratePFK(password string, domain string) {
	cmd := []string{"openssl", "pkcs12", "-export", "-out", domain + ".pfx", "-inkey", domain + ".key", "-in", domain + ".pem", "-passin", "pass:" + password + "", "-passout", "pass:" + password + ""}
	err := Cmd(strings.Join(cmd, " "))
	if err != nil {
		logrus.Info("[-] cmd.Run() failed with", err.Error())
		os.Exit(0)
	}
}

func SignExecutable(domain string, filein string) {
	if domain == "" {
		return
	}

	if _, err := exec.LookPath("openssl"); err != nil {
		logrus.Warn("[-] 缺少openssl命令")
		return
	}

	if _, err := exec.LookPath("osslsigncode"); err != nil {
		logrus.Warn("[-] 缺少osslsigncode签名")
		return
	}

	os.Rename(filein, filein+".old")
	inputFile := filein + ".old"
	defer os.Remove(inputFile)

	pfx := domain
	password := strconv.FormatInt(time.Now().Unix(), 10)

	GenerateCert(domain)
	GeneratePFK(password, domain)
	pfx = domain + ".pfx"
	logrus.Warn("[+] 生成签名证书:", pfx)

	logrus.Info("[+] 签名文件" + filein)
	cmd := []string{"osslsigncode", "sign", "-pkcs12", pfx, "-in", "" + inputFile + "", "-out", "" + filein + "", "-pass", "" + password + ""}
	err := Cmd(strings.Join(cmd, " "))
	if err != nil {
		logrus.Info("[-] cmd.Run() failed with", err.Error())
		os.Exit(0)
	}
}
