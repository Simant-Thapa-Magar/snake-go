package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

type Coordinate struct {
	x, y int
}

type Snake struct {
	points                      []*Coordinate
	columnVelocity, rowVelocity int
	symbol                      rune
}

type Apple struct {
	point  *Coordinate
	symbol rune
}

var snake *Snake
var apple *Apple
var Screen tcell.Screen
var screenWidth, screenHeight int

const FRAME_WIDTH = 80
const FRAME_HEIGHT = 15
const FRAME_BORDER_THICKNESS = 1
const FRAME_BORDER_VERTICAL = '║'
const FRAME_BORDER_HORIZONTAL = '═'
const FRAME_BORDER_TOP_LEFT = '╔'
const FRAME_BORDER_TOP_RIGHT = '╗'
const FRAME_BORDER_BOTTOM_RIGHT = '╝'
const FRAME_BORDER_BOTTOM_LEFT = '╚'
const SNAKE_SYMBOL = 0x2588
const APPLE_SYMBOL = 0x25CF

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	initScreen()
	initializeGameObjects()
	userInput := readUserInput()
	var key string
	for {
		displayFrame()
		key = getUserInput(userInput)
		handleUserInput(key)
		updateGameState()
		displayGameObjects()
		time.Sleep(75 * time.Millisecond)
		Screen.Clear()
	}
}

func getUserInput(userInput chan string) string {
	var key string
	select {
	case key = <-userInput:
	default:
		key = ""
	}
	return key
}

func readUserInput() chan string {
	userInput := make(chan string)
	go func() {
		for {
			switch ev := Screen.PollEvent().(type) {
			case *tcell.EventKey:
				userInput <- ev.Name()
			}
		}
	}()
	return userInput
}

func handleUserInput(key string) {
	if key == "Rune[q]" {
		Screen.Fini()
		os.Exit(0)
	} else if key == "Up" && snake.rowVelocity == 0 {
		snake.rowVelocity = -1
		snake.columnVelocity = 0
	} else if key == "Down" && snake.rowVelocity == 0 {
		snake.rowVelocity = 1
		snake.columnVelocity = 0
	} else if key == "Left" && snake.columnVelocity == 0 {
		snake.rowVelocity = 0
		snake.columnVelocity = -1
	} else if key == "Right" && snake.columnVelocity == 0 {
		snake.rowVelocity = 0
		snake.columnVelocity = 1
	}
}

func getSnakeHeadCoordinates() (int, int) {
	snakeHead := snake.points[len(snake.points)-1]
	return snakeHead.x, snakeHead.y
}

func updateSnake() {
	snakeHeadX, snakeHeadY := getSnakeHeadCoordinates()
	snake.points = append(snake.points, &Coordinate{
		snakeHeadX + snake.columnVelocity,
		snakeHeadY + snake.rowVelocity,
	})
	if !isAppleInsideSnake() {
		snake.points = snake.points[1:]
	}
	updateSnakeIfBeyoundBorder()
}

func updateSnakeIfBeyoundBorder() {
	originX, originY := getFrameOrigin()
	topY := originY
	bottomY := originY + FRAME_HEIGHT
	leftX := originX
	rightX := originX + FRAME_WIDTH - 1
	for _, snakeCoordinate := range snake.points {
		if snakeCoordinate.y <= topY {
			// if above
			snakeCoordinate.y = bottomY - 1
		} else if snakeCoordinate.y >= bottomY {
			// if below
			snakeCoordinate.y = topY + 1
		} else if snakeCoordinate.x >= rightX {
			// if right
			snakeCoordinate.x = leftX + 1
		} else if snakeCoordinate.x <= leftX {
			// if left
			snakeCoordinate.x = rightX - 1
		}
	}
}

func isAppleInsideSnake() bool {
	for _, snakeCoordinate := range snake.points {
		if snakeCoordinate.x == apple.point.x && snakeCoordinate.y == apple.point.y {
			return true
		}
	}
	return false
}

func getNewAppleCoordinate() (int, int) {
	rand.Seed(time.Now().UnixMicro())
	randomX := rand.Intn(FRAME_WIDTH)
	randomY := rand.Intn(FRAME_HEIGHT)

	newCoordinate := &Coordinate{
		randomX, randomY,
	}

	transformCoordinateInsideFrame(newCoordinate)

	return newCoordinate.x, newCoordinate.y
}

func updateApple() {
	for isAppleInsideSnake() {
		apple.point.x, apple.point.y = getNewAppleCoordinate()
	}
}

func updateGameState() {
	updateSnake()
	updateApple()
}

func transformCoordinateInsideFrame(coordinate *Coordinate) {
	frameOriginX, frameOriginY := getFrameOrigin()
	coordinate.x += frameOriginX
	coordinate.y += frameOriginY
}

func initializeGameObjects() {
	snake = &Snake{
		points:         getInitialSnakeCoordinates(),
		columnVelocity: 0,
		rowVelocity:    1,
		symbol:         SNAKE_SYMBOL,
	}

	apple = &Apple{
		point:  getInitialAppleCoordinates(),
		symbol: APPLE_SYMBOL,
	}
}

func getInitialSnakeCoordinates() []*Coordinate {
	snakeInitialCoordinate1 := &Coordinate{8, 4}
	transformCoordinateInsideFrame(snakeInitialCoordinate1)

	snakeInitialCoordinate2 := &Coordinate{8, 5}
	transformCoordinateInsideFrame(snakeInitialCoordinate2)

	snakeInitialCoordinate3 := &Coordinate{8, 6}
	transformCoordinateInsideFrame(snakeInitialCoordinate3)

	snakeInitialCoordinate4 := &Coordinate{8, 7}
	transformCoordinateInsideFrame(snakeInitialCoordinate4)

	return []*Coordinate{
		{snakeInitialCoordinate1.x, snakeInitialCoordinate1.y},
		{snakeInitialCoordinate2.x, snakeInitialCoordinate2.y},
		{snakeInitialCoordinate3.x, snakeInitialCoordinate3.y},
		{snakeInitialCoordinate4.x, snakeInitialCoordinate4.y},
	}
}

func getInitialAppleCoordinates() *Coordinate {
	appleInitialCoordinate := &Coordinate{FRAME_WIDTH / 2, FRAME_HEIGHT / 2}
	transformCoordinateInsideFrame(appleInitialCoordinate)

	return appleInitialCoordinate
}

func initScreen() {
	encoding.Register()
	var err error
	Screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = Screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	Screen.SetStyle(defStyle)
	screenWidth, screenHeight = Screen.Size()
}

func print(x, y, w, h int, char rune) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			Screen.SetContent(x+i, y+j, char, nil, tcell.StyleDefault)
		}
	}
}

// func displayHelloWorld(s tcell.Screen) {
// 	w, h := s.Size()
// 	s.Clear()
// 	style := tcell.StyleDefault.Foreground(tcell.ColorCadetBlue.TrueColor()).Background(tcell.ColorWhite)
// 	emitStr(s, w/2-7, h/2, style, "Hello, World!")
// 	emitStr(s, w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")
// 	s.Show()
// }

func getFrameOrigin() (int, int) {
	return (screenWidth-FRAME_WIDTH)/2 - 1, (screenHeight-FRAME_HEIGHT)/2 - 1
}

func displayFrame() {
	frameOriginX, frameOriginY := getFrameOrigin()
	printUnfilledRectangle(frameOriginX, frameOriginY, FRAME_WIDTH, FRAME_HEIGHT, FRAME_BORDER_THICKNESS, FRAME_BORDER_HORIZONTAL, FRAME_BORDER_VERTICAL, FRAME_BORDER_TOP_LEFT, FRAME_BORDER_TOP_RIGHT, FRAME_BORDER_BOTTOM_RIGHT, FRAME_BORDER_BOTTOM_LEFT)
	Screen.Show()
}

func displayGameObjects() {
	displaySnake()
	displayApple()
	Screen.Show()
}

func displaySnake() {
	for _, snakeCoordinate := range snake.points {
		print(snakeCoordinate.x, snakeCoordinate.y, 1, 1, snake.symbol)
	}
}

func displayApple() {
	print(apple.point.x, apple.point.y, 1, 1, apple.symbol)
}

func printUnfilledRectangle(xOrigin, yOrigin, width, height, borderThickness int, horizontalOutline, verticalOutline, topLeftOutline, topRightOutline, bottomRightOutline, bottomLeftOutline rune) {
	var upperBorder, lowerBorder rune
	verticalBorder := verticalOutline
	for i := 0; i < width; i++ {
		// upper boundry
		if i == 0 {
			upperBorder = topLeftOutline
			lowerBorder = bottomLeftOutline
		} else if i == width-1 {
			upperBorder = topRightOutline
			lowerBorder = bottomRightOutline
		} else {
			upperBorder = horizontalOutline
			lowerBorder = horizontalOutline
		}
		print(xOrigin+i, yOrigin, borderThickness, borderThickness, upperBorder)
		print(xOrigin+i, yOrigin+height, borderThickness, borderThickness, lowerBorder)
		// lower boundry
	}

	// side boundry
	for i := 1; i < height; i++ {
		print(xOrigin, yOrigin+i, borderThickness, borderThickness, verticalBorder)
		print(xOrigin+width-1, yOrigin+i, borderThickness, borderThickness, verticalBorder)
	}
}
