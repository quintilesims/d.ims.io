package controllers

import (
	"log"

	"github.com/zpatrick/fireball"
)

func LogDecorator() fireball.Decorator {
	return func(handler fireball.Handler) fireball.Handler {
		return func(c *fireball.Context) (fireball.Response, error) {
			log.Printf("[DEBUG] %s %s %s\n",
				c.Request.RemoteAddr,
				c.Request.Method,
				c.Request.URL.String())

			return handler(c)
		}
	}
}
