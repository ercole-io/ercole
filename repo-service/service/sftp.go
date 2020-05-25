// Copyright (c) 2019 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package service is a package that contains varios file serving services
package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"

	"github.com/ercole-io/ercole/config"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// SFTPRepoSubService is a concrete implementation of SubRepoServiceInterface
type SFTPRepoSubService struct {
	// Config contains the reposervice global configuration
	Config config.Configuration
	// Log contains hs.Logger formatted
	Log *logrus.Logger
}

// Init start the service
func (hs *SFTPRepoSubService) Init(_ *sync.WaitGroup) {
	//FIXME: Temporary fix
	if err := os.Chdir("/"); err != nil {
		hs.Log.Fatal("Cannot change directory to /", err)
	}
	//Setup the ssh server config
	privateKeyBytes, err := ioutil.ReadFile(hs.Config.RepoService.SFTP.PrivateKey)
	if err != nil {
		hs.Log.Fatal("Failed to load the repo-service/ssh private key", err)
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		hs.Log.Fatal("Failed to parse the repo-service/ssh private key", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(privateKey)

	//start the listener
	hs.Log.Info("Start repo-service/sftp: listening at ", hs.Config.RepoService.SFTP.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", hs.Config.RepoService.SFTP.BindIP, hs.Config.RepoService.SFTP.Port))
	if err != nil {
		hs.Log.Fatal("Stopping repo-service/http", err)
	}

	//Start the sftp sub service
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			hs.Log.Infof("Failed to accept incoming connection (%s)", err)
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			hs.Log.Infof("Failed to handshake (%s)", err)
			continue
		}

		if hs.Config.RepoService.SFTP.LogConnections {
			hs.Log.Infof("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		}
		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)
		// Accept all channels
		go hs.handleChannels(chans)
	}
}

// handleChannels will handle any incoming new channels
func (hs *SFTPRepoSubService) handleChannels(chans <-chan ssh.NewChannel) {
	//serve the incoming Channel channel
	for newChannel := range chans {
		go hs.handleChannel(newChannel)
	}
}

// handleChannel handle a single channel
func (hs *SFTPRepoSubService) handleChannel(newChannel ssh.NewChannel) {

	//Check the channel
	if hs.Config.RepoService.SFTP.DebugConnections {
		hs.Log.Infof("Incoming channel: %s\n", newChannel.ChannelType())
	}
	if newChannel.ChannelType() != "session" {
		newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
		hs.Log.Infof("Unknown channel type: %s\n", newChannel.ChannelType())
		return
	}

	//Accept the channel
	channel, requests, err := newChannel.Accept()
	if err != nil {
		hs.Log.Fatal("could not accept channel.", err)
	}
	if hs.Config.RepoService.SFTP.DebugConnections {
		hs.Log.Info("Channel accepted\n")
	}

	//Handle the request
	go func(in <-chan *ssh.Request) {
		for req := range in {
			if hs.Config.RepoService.SFTP.DebugConnections {
				hs.Log.Infof("Request: %v\n", req.Type)
			}
			ok := false
			switch req.Type {
			case "subsystem":
				if hs.Config.RepoService.SFTP.DebugConnections {
					hs.Log.Infof("Subsystem: %s\n", req.Payload[4:])
				}
				if string(req.Payload[4:]) == "sftp" {
					ok = true
				}
			}
			if hs.Config.RepoService.SFTP.DebugConnections {
				hs.Log.Infof(" - accepted: %v\n", ok)
			}
			req.Reply(ok, nil)
		}
	}(requests)

	//Setup the the sftp server
	serverOptions := []sftp.ServerOption{
		sftp.ReadOnly(),
		sftp.RootDirectory(hs.Config.RepoService.DistributedFiles),
	}
	if hs.Config.RepoService.SFTP.DebugConnections {
		serverOptions = append(serverOptions, sftp.WithDebug(os.Stdout))
	}
	server, err := sftp.NewServer(
		channel,
		serverOptions...,
	)
	if err != nil {
		hs.Log.Fatal(err)
	}

	//Serve the sftp client
	if err := server.Serve(); err == io.EOF {
		server.Close()
		if hs.Config.RepoService.SFTP.LogConnections {
			hs.Log.Info("sftp client exited session.")
		}
	} else if err != nil {
		hs.Log.Fatal("sftp server completed with error:", err)
	}

}
