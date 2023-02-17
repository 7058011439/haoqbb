package server

import (
	"github.com/7058011439/haoqbb/Ftp/server/core"
	"github.com/7058011439/haoqbb/Ftp/server/driver/file"
)

type (
	Auth                  = core.Auth
	Command               = core.Command
	Conn                  = core.Conn
	DataSocket            = core.DataSocket
	DiscardLogger         = core.DiscardLogger
	Driver                = core.Driver
	DriverFactory         = core.DriverFactory
	FileInfo              = core.FileInfo
	Logger                = core.Logger
	MultipleDriver        = core.MultipleDriver
	MultipleDriverFactory = core.MultipleDriverFactory
	Notifier              = core.Notifier
	NullNotifier          = core.NullNotifier
	Perm                  = core.Perm
	Server                = core.Server
	Opts                  = core.ServerOpts
	SimpleAuth            = core.SimpleAuth
	SimplePerm            = core.SimplePerm
	StdLogger             = core.StdLogger
)

var (
	Version = core.Version
)

type (
	FileDriver        = file.Driver
	FileDriverFactory = file.DriverFactory
)
