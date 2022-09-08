package http

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/valyala/fasthttp"
	"github.com/xxxmailk/cera/log"
	"github.com/xxxmailk/cera/middlewares"
	"github.com/xxxmailk/cera/router"
	"math/big"
	"net"
	"os"
	"time"
)

type StartHttpServer interface {
	SetLogger(l log.SimpleLogger)
	SetHostname(hostname string)
	SetRouter(handler *router.Router)
	SetIdleTimeout(sec int)
	Start() error
	Stop() error
	UseMiddleWare(m middlewares.MiddlewareInterface)
	AtLast(m middlewares.MiddlewareInterface)
}

type StartTlsServer interface {
	SetLogger(l log.SimpleLogger)
	SetHostname(hostname string)
	SetRouter(handler *router.Router)
	SetIdleTimeout(sec int)
	StartTls() error
	Stop() error
	AtLast(m middlewares.MiddlewareInterface)
	UseMiddleWare(m middlewares.MiddlewareInterface)
}

type Serve struct {
	ip          string
	port        string
	idleTimeout time.Duration
	hostname    string
	logger      log.SimpleLogger
	handler     fasthttp.RequestHandler
	sslKey      string
	sslCert     string
	router      *router.Router
	middleWares []middlewares.MiddlewareInterface
	lastFunc    []middlewares.MiddlewareInterface
	serv        *fasthttp.Server
}

func (s *Serve) SetIdleTimeout(sec int) {
	s.idleTimeout = time.Duration(sec) * time.Second
}

func (s *Serve) SetHostname(hostname string) {
	s.hostname = hostname
}

func (s *Serve) SetHandle(h fasthttp.RequestHandler) {
	s.handler = h
}

func (s *Serve) SetLogger(l log.SimpleLogger) {
	s.logger = l
}
func (s *Serve) Stop() error {
	return s.serv.Shutdown()
}

func (s *Serve) SetSslKeyCert(keyPath, certPath string) {
	s.sslKey, s.sslCert = keyPath, certPath
}

func (s *Serve) SetRouter(handler *router.Router) {
	s.router = handler
}

func (s *Serve) Start() error {
	s.SetHandle(s.httpHandler)
	s.serv = &fasthttp.Server{
		// allocation http handle with domain name
		Handler:     s.handler,
		IdleTimeout: s.idleTimeout,
	}
	if s.router == nil {
		panic("please set router before server start server")
	}
	s.router.Logger = s.logger
	s.logger.Infof("starting web server and listening on %s:%s", s.ip, s.port)
	return s.ListenAndServe()
}

func (s *Serve) ListenAndServe() error {
	addr := net.JoinHostPort(s.ip, s.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.serv.Serve(ln)
}

func (s *Serve) StartTls() error {
	s.SetHandle(s.httpHandler)
	server := &fasthttp.Server{
		// allocation http handle with domain name
		Handler:     s.handler,
		IdleTimeout: s.idleTimeout,
	}

	if s.sslCert == "" && s.sslKey == "" {

		// preparing second host
		cert, priv, err := GenerateCert(net.JoinHostPort("127.0.0.1", s.port))

		if s.hostname != "" {
			cert, priv, err = GenerateCert(s.hostname)
		} else {
			// todo: hostname config
			s.logger.Warnf("hostname has not been configured, use 127.0.0.1:%s to generate ssl certs", s.port)
		}

		if err != nil {
			s.logger.Errorf(err.Error())
			return err
		}

		err = server.AppendCertEmbed(cert, priv)
		if err != nil {
			s.logger.Errorf(err.Error())
			return err
		}
		if s.router == nil {
			panic("please set router before server start server")
		}
		s.logger.Infof("starting TLS web server and listening on %s:%s", s.ip, s.port)
		err = server.ListenAndServeTLS(s.port, "", "")
		if err != nil {
			s.logger.Errorf(err.Error())
		}
		return err
	}
	s.logger.Infof("starting TLS web server and listening on %s:%s", s.ip, s.port)

	// if ssl cert and ssl key had been set, use cert and key file to start ssl server
	err := server.ListenAndServeTLS(s.port, s.sslCert, s.sslKey)
	if err != nil {
		s.logger.Fatalf(err.Error())
		return err
	}
	return nil
}

// new simple http server
// you can set your http server before start()
func NewHttpServe(ip, port string) StartHttpServer {
	l := log.NewSimpleLogger()
	host, err := os.Hostname()
	if err != nil {
		l.Fatalf(err.Error())
	}
	s := &Serve{
		ip:          ip,
		port:        port,
		hostname:    host,
		idleTimeout: time.Duration(30) * time.Second,
		logger:      l,
	}
	return s
}

// ditto ↑
func NewTLSServe(ip, port string) StartTlsServer {
	l := log.NewSimpleLogger()
	host, err := os.Hostname()
	if err != nil {
		l.Fatalf(err.Error())
	}
	s := &Serve{
		ip:          ip,
		port:        port,
		hostname:    host,
		idleTimeout: time.Duration(30) * time.Second,
		logger:      l,
	}
	return s
}

func GenerateCert(host string) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"CERA"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		DNSNames:              []string{host},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader, cert, cert, &priv.PublicKey, priv,
	)

	p := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	b := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		},
	)

	return b, p, err
}

func (s *Serve) UseMiddleWare(m middlewares.MiddlewareInterface) {
	s.middleWares = append(s.middleWares, m)
}

func (s *Serve) AtLast(m middlewares.MiddlewareInterface) {
	s.lastFunc = append(s.lastFunc, m)
}

func (s *Serve) httpHandler(ctx *fasthttp.RequestCtx) {
	if err := recover(); err != nil {
		panic(err)
	}

	// handling middleWares
	if len(s.middleWares) > 0 {
		for i := len(s.middleWares); i >= 0; i-- {
			if i > 0 {
				mid := s.middleWares[i-1]
				ctx = mid.Handle(ctx)
				if mid.IsBreakHere() {
					return
				}
			}
		}
	}
	s.logger.Debugf("handle with router handler")
	// transfer http contexts to router handler
	s.router.Handler(ctx)

	// handling last functions
	if len(s.lastFunc) > 0 {
		for i := len(s.lastFunc); i >= 0; i-- {
			if i > 0 {
				last := s.lastFunc[i-1]
				last.Handle(ctx)
				if last.IsBreakHere() {
					return
				}
			}
		}
	}
}
