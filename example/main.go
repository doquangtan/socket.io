package main

import (
	"log"

	socketio "github.com/doquangtan/gofiber-socket.io/v4"
	"github.com/gofiber/fiber/v2"
)

func socketIoRoute(app fiber.Router) {
	io := socketio.New()

	io.OnAuthorization(func(params map[string]string) bool {
		// auth, ok := params["Authorization"]
		// if !ok {
		// 	return false
		// }
		return true
	})

	io.OnConnection(func(socket *socketio.Socket) {
		println("connect", socket.Nps, socket.Id)
		socket.Join("demo")
		io.To("demo").Emit("test", socket.Id+" join us room...", "server message")

		socket.On("connected", func(event *socketio.EventPayload) {
			socket.Emit("chat message", "Main")
		})
		socket.On("test", func(event *socketio.EventPayload) {
			socket.Emit("test", event.Data...)
		})

		socket.On("join-room", func(event *socketio.EventPayload) {
			if len(event.Data) > 0 && event.Data[0] != nil {
				socket.Join(event.Data[0].(string))
			}
		})

		socket.On("to-room", func(event *socketio.EventPayload) {
			socket.To("demo").To("demo2").Emit("test", "hello")
		})

		socket.On("leave-room", func(event *socketio.EventPayload) {
			socket.Leave("demo")
			socket.Join("demo2")
		})

		socket.On("my-room", func(event *socketio.EventPayload) {
			socket.Emit("my-room", socket.Rooms())
		})

		socket.On("chat message", func(event *socketio.EventPayload) {
			socket.Emit("chat message", event.Data[0])

			if len(event.Data) > 2 {
				log.Println(socket.Nps, ": ", event.Data[2].(map[string]interface{}))
			}

			if event.Ack != nil {
				event.Ack("hello from name space root", map[string]interface{}{
					"Test": "ok",
				})
			}
		})

		socket.On("disconnecting", func(event *socketio.EventPayload) {
			println("disconnecting", socket.Nps, socket.Id)
		})

		socket.On("disconnect", func(event *socketio.EventPayload) {
			println("disconnect", socket.Nps, socket.Id)
		})
	})

	io.Of("/test").OnConnection(func(socket *socketio.Socket) {
		println("connect", socket.Nps, socket.Id)

		socket.On("chat message", func(event *socketio.EventPayload) {
			socket.Emit("chat message", event.Data[0])

			if len(event.Data) > 2 {
				log.Println(socket.Nps, ": ", event.Data[2].(map[string]interface{}))
			}

			if event.Ack != nil {
				event.Ack("hello from nps test", map[string]interface{}{
					"Test": "ok",
				})
			}
		})
	})

	app.Use("/", io.Middleware)
	app.Route("/socket.io", io.Server)
}

func main() {
	app := fiber.New(fiber.Config{})

	app.Static("/", "./public")

	app.Route("/", socketIoRoute)

	app.Get("/test", func(c *fiber.Ctx) error {
		io := c.Locals("io").(*socketio.Io)

		io.Emit("event", map[string]interface{}{
			"Ok": 1,
		})

		io.Of("/admin").Emit("event", map[string]interface{}{
			"Ok": 1,
		})

		return c.SendStatus(200)
	})

	app.Listen(":3300")
}
