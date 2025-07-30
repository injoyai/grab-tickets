package node

import (
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/server"
	"github.com/injoyai/ios/server/listen"
)

func Listen(port int) error {
	return listen.RunTCP(port, func(s *server.Server) {
		s.SetClientOption(func(c *client.Client) {
			c.OnDealMessage = func(c *client.Client, msg ios.Acker) {

			}
		})
	})
}
