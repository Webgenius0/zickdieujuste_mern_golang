package main

import "github.com/labstack/echo/v5"

func main() {
    e := echo.New()
    e.GET("/", func(c *echo.Context) error {
        return c.JSON(200, map[string]any{"message": "Hello, World!"})
    })
    e.Start(":1323")
}