package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ckoch786/fiber/v3"
)

// TODO what does this buy use?
// Why do so many types contain a pointer to the App?
type server struct {
	app *fiber.App
}

var s server

func Runner() {
	// Prompt user to select which option to run of the two function calls below
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Which option would you like to run?")
	fmt.Println("1. Allow Removing Registered Routing")
	fmt.Println("2. Repo Example")

	for {
		fmt.Print("Enter option number: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		// Sanitize input
		input = strings.TrimSpace(input)

		if input == "1" {
			AllowRemovingRegisteredRouting()
			break
		} else if input == "2" {
			RepoExample()
			break
		} else {
			fmt.Println("Invalid option. Please enter 1 or 2.")
		}
	}
}

func AllowRemovingRegisteredRouting() {
	s.app = fiber.New()
	s.app.Get("/define", func(c fiber.Ctx) error {
		before := s.app.HandlersCount()
		addMyRoute()
		after := s.app.HandlersCount()

		return c.JSON(map[string]any{
			"hcBefore": before,
			"hcAfter":  after,
		})
	})

	log.Fatal(s.app.Listen(":3000"))
}

func addMyRoute() {
	// Remove route if it already exists
	// TODO why would this and AddRoute need to aquire a lock in order to be performed? Is there a case
	// where this would be called in such a way that would cause a deadlock?  What happens to the blocked calls
	// that are waiting to aquire a lock on the app.mutex?
	s.app.RemoveRoute("/hello", fiber.MethodHead, fiber.MethodGet)
	// TODO In this case is there a way to check to see if the route is already in the
	// stack and just not add it again? Or should the remove method do nothing if the
	// route does not already exist.

	// TODO look into tests around all methods that are called for clues.
	// Add new route
	s.app.Get("/hello", func(c fiber.Ctx) error {
		return c.SendString("hello")
	})

	// TODO comment this out and see how the server behaves
	// Rebuild the routing tree to reflect changes
	// TOODO what is resource intensive about this?
	s.app.RebuildTree() // Seems to do nothing?
	// ^ This rebuild tree is needed because the build function is only called by the *app.startupProcess
	// since this route is dynamically defined after the server has been started.
}

func main() {
	AllowRemovingRegisteredRouting()
	// Runner()
}

func RepoExample() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
