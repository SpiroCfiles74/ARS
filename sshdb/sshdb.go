package SSHdb

import (
	"bytes"
	"crypto/x509"
	//"database/sql"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

// Declare your connection data and user credentials here
const (
	// ssh connection related data
	sshServerHost     = "208.75.75.70"
	sshServerPort     = 22
	sshUserName       = "ubuntu"
	sshPrivateKeyFile = "ars_key_pair.pem" // exported as OpenSSH key from .ppk
	sshKeyPassphrase  = "Test12321"        // key file encrytion password

	// ssh tunneling related data
	sshLocalHost  = "localhost" // local localhost ip (client side)
	sshLocalPort  = 9000        // local port used to forward the connection
	sshRemoteHost = "127.0.0.1" // remote local ip (server side)
	sshRemotePort = 3306        // remote MySQL port

	// MySQL access data
	mySqlUsername = "root"
	mySqlPassword = "Test12321"
	mySqlDatabase = "arsystems"
)

//dbErrorHandler - Simple mySql error handling (yet to implement)
func dbErrorHandler(err error) {
	switch err := err.(type) {
	default:
		fmt.Printf("Error %s\n", err)
		os.Exit(-1)
	}
}

// Endpoint - Define an endpoint with ip and port
type Endpoint struct {
	Host string
	Port int
}

// Returns an endpoint as ip:port formatted string
func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// SSHtunnel - Define the endpoints along the tunnel
type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
}

// Start the tunnel
func (tunnel *SSHtunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

// Port forwarding
func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	// Establish connection to the intermediate server
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	// access the target server
	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	// Transfer the data between  and the remote server
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

// DecryptPEMkey : Decrypt encrypted PEM key data with a passphrase and embed it to key prefix
// and postfix header data to make it valid for further private key parsing.
func DecryptPEMkey(buffer []byte, passphrase string) []byte {
	block, _ := pem.Decode(buffer)
	der, err := x509.DecryptPEMBlock(block, []byte(passphrase))
	if err != nil {
		fmt.Println("decrypt failed: ", err)
	}
	encoded := base64.StdEncoding.EncodeToString(der)
	encoded = "-----BEGIN RSA PRIVATE KEY-----\n" + encoded +
		"\n-----END RSA PRIVATE KEY-----\n"
	return []byte(encoded)
}

// PublicKeyFile - Get the signers from the OpenSSH key file (.pem) and return them for use in
// the Authentication method. Decrypt encrypted key data with the passphrase.*/
func PublicKeyFile(file string, passphrase string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	if bytes.Contains(buffer, []byte("ENCRYPTED")) {
		// Decrypt the key with the passphrase if it has been encrypted
		buffer = DecryptPEMkey(buffer, passphrase)
	}

	// Get the signers from the key
	signers, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(signers)
}

// Define the ssh tunnel using its endpoint and config data
func sshTunnel() *SSHtunnel {
	localEndpoint := &Endpoint{
		Host: sshLocalHost,
		Port: sshLocalPort,
	}

	serverEndpoint := &Endpoint{
		Host: sshServerHost,
		Port: sshServerPort,
	}

	remoteEndpoint := &Endpoint{
		Host: sshRemoteHost,
		Port: sshRemotePort,
	}

	sshConfig := &ssh.ClientConfig{
		User: sshUserName,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(sshPrivateKeyFile, sshKeyPassphrase)},
	}

	return &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}
}
