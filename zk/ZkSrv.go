package zk

import (
	"time"
	"github.com/samuel/go-zookeeper/zk"
	_ "github.com/nikofake/im-gateway/yml"
	"github.com/nikofake/im-gateway/yml"
	"github.com/gogap/logrus"
	"unicode/utf8"
	"strings"
	"github.com/nikofake/im-gateway/ip"
	"encoding/json"
	"errors"
)

type HostInfo struct {
	Ip   string
	Port int32
}

func init() {
	if c, _, err := zk.Connect(yml.Config.Zk.Hosts, time.Second); err != nil {
		panic(err)
	} else {
		monitorPath := yml.Config.Zk.Basepath + "/monitor"
		if _, err := cretePathRecur(c, monitorPath, 0); err != nil {
			panic(err)
		}

		connectionPath := yml.Config.Zk.Basepath + "/connection"
		if _, err := cretePathRecur(c, connectionPath, 0); err != nil {
			panic(err)
		}

		monitorBytes, _ := json.Marshal(HostInfo{Ip: ip.GetInternal(), Port: yml.Config.Monitor.Port})
		if path, error := c.Create(monitorPath+"/"+string(monitorBytes), monitorBytes, zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); error != nil {
			panic(error)
		} else {
			logrus.Info("reg monitor srv to zk ", path)
		}

		connectionBytes, _ := json.Marshal(HostInfo{Ip: ip.GetInternal(), Port: yml.Config.Connection.Port})
		if path, error := c.Create(connectionPath+"/"+string(connectionBytes), connectionBytes, zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); error != nil {
			panic(error)
		} else {
			logrus.Info("reg connection srv to zk", path)
		}
	}

}

func cretePathRecur(c *zk.Conn, path string, flag int32) (p string, error error) {

	if validatePath(path) {
		temp := ""
		for _, e := range strings.Split(path, "/") {
			if len(e) > 0 {
				temp = temp + "/" + e
				if exists, _, err := c.Exists(temp); err != nil {
					error = err
					return
				} else {
					if !exists {
						if _, err := c.Create(temp, nil, flag, zk.WorldACL(zk.PermAll)); err != nil {
							error = err
							return
						}
					}
				}
			}
		}
		p = temp
	} else {
		error = errors.New("invalid path " + path)
	}

	return
}

func validatePath(path string) bool {
	if path == "" {
		return false
	}

	if path[0] != '/' {
		return false
	}

	n := len(path)
	if n == 1 {
		// path is just the root
		return false
	}

	if path[n-1] == '/' {
		return false
	}

	// Start at rune 1 since we already know that the first character is
	// a '/'.
	for i, w := 1, 0; i < n; i += w {
		r, width := utf8.DecodeRuneInString(path[i:])
		switch {
		case r == '\u0000':
			return false
		case r == '/':
			last, _ := utf8.DecodeLastRuneInString(path[:i])
			if last == '/' {
				return false
			}
		case r == '.':
			last, lastWidth := utf8.DecodeLastRuneInString(path[:i])

			// Check for double dot
			if last == '.' {
				last, _ = utf8.DecodeLastRuneInString(path[:i-lastWidth])
			}

			if last == '/' {
				if i+1 == n {
					return false
				}

				next, _ := utf8.DecodeRuneInString(path[i+w:])
				if next == '/' {
					return false
				}
			}
		case r >= '\u0000' && r <= '\u001f',
			r >= '\u007f' && r <= '\u009f',
			r >= '\uf000' && r <= '\uf8ff',
			r >= '\ufff0' && r < '\uffff':
			return false
		}
		w = width
	}
	return true
}
