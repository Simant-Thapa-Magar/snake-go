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
var coordinatesToClear []*Coordinate
var Screen tcell.Screen
var screenWidth, screenHeight int
var isGamePaused, isGameOver bool
var score int

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
	displayFrame()
	displayGameScore()
	userInput := readUserInput()
	var key string
	for !isGameOver {
		if isGamePaused {
			displayGamePausedInfo()
		}
		key = getUserInput(userInput)
		handleUserInput(key)
		updateGameState()
		displayGameObjects()
		time.Sleep(75 * time.Millisecond)
	}

	displayGameOverInfo()
	time.Sleep(3 * time.Second)
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
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	} else if !isGamePaused {
		if key == "Up" && snake.rowVelocity == 0 {
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
}

func getSnakeHeadCoordinates() (int, int) {
	snakeHead := snake.points[len(snake.points)-1]
	return snakeHead.x, snakeHead.y
}

func updateSnake() {
	snakeHeadX, snakeHeadY := getSnakeHeadCoordinates()
	newSnakeHead := &Coordinate{
		snakeHeadX + snake.columnVelocity,
		snakeHeadY + snake.rowVelocity,
	}
	setSnakeWithinFrame(newSnakeHead)
	snake.points = append(snake.points, newSnakeHead)
	if !isAppleInsideSnake() {
		coordinatesToClear = append(coordinatesToClear, snake.points[0])
		snake.points = snake.points[1:]
	} else {
		score++
		displayGameScore()
	}
	if isSnakeEatingItself() {
		isGameOver = true
	}
}

func setSnakeWithinFrame(snakeCoordinate *Coordinate) {
	originX, originY := getFrameOrigin()
	topY := originY
	bottomY := originY + FRAME_HEIGHT - 1
	leftX := originX
	rightX := originX + FRAME_WIDTH - 1

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

func isSnakeEatingItself() bool {
	headX, headY := getSnakeHeadCoordinates()
	for _, snakeCoordinate := range snake.points[:len(snake.points)-1] {
		if headX == snakeCoordinate.x && headY == snakeCoordinate.y {
			return true
		}
	}
	return false
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
	randomX := rand.Intn(FRAME_WIDTH - 2)
	randomY := rand.Intn(FRAME_HEIGHT - 2)

	newCoordinate := &Coordinate{
		randomX, randomY,
	}

	transformCoordinateInsideFrame(newCoordinate)

	return newCoordinate.x, newCoordinate.y
}

func updateApple() {
	for isAppleInsideSnake() {
		coordinatesToClear = append(coordinatesToClear, apple.point)
		apple.point.x, apple.point.y = getNewAppleCoordinate()
	}
}

func updateGameState() {
	if isGamePaused {
		return
	}
	clearScreen()
	updateSnake()
	updateApple()
}

func transformCoordinateInsideFrame(coordinate *Coordinate) {
	frameOriginX, frameOriginY := getFrameOrigin()
	frameOriginX += 1
	frameOriginY += 1
	coordinate.x += frameOriginX
	coordinate.y += frameOriginY
	for coordinate.x >= frameOriginX+FRAME_WIDTH {
		coordinate.x--
	}
	for coordinate.y >= frameOriginY+FRAME_HEIGHT-1 {
		coordinate.y--
	}
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

	if screenWidth < FRAME_WIDTH || screenHeight < FRAME_HEIGHT {
		fmt.Printf("The game frame is defined with %d width and %d height. Increase terminal size and try again ", FRAME_WIDTH, FRAME_HEIGHT)
		os.Exit(1)
	}
}

func print(x, y, w, h int, style tcell.Style, char rune) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			Screen.SetContent(x+i, y+j, char, nil, style)
		}
	}
}

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
	style := tcell.StyleDefault.Foreground(tcell.ColorDarkGreen.TrueColor())
	for _, snakeCoordinate := range snake.points {
		print(snakeCoordinate.x, snakeCoordinate.y, 1, 1, style, snake.symbol)
	}
}

func displayApple() {
	style := tcell.StyleDefault.Foreground(tcell.ColorDarkRed.TrueColor())
	print(apple.point.x, apple.point.y, 1, 1, style, apple.symbol)
}

func displayGamePausedInfo() {
	_, frameY := getFrameOrigin()
	printAtCenter(frameY-2, "Game Paused !!", true)
	printAtCenter(frameY-1, "Press p to resume", true)
}

func displayGameOverInfo() {
	centerY := (screenHeight - FRAME_HEIGHT) / 2
	printAtCenter(centerY-1, "Game Over !!", false)
	printAtCenter(centerY, fmt.Sprintf("Your Score : %d", score), false)
}

func displayGameScore() {
	_, frameY := getFrameOrigin()
	printAtCenter(frameY+FRAME_HEIGHT+2, fmt.Sprintf("Current Score : %d", score), false)
}

func printAtCenter(startY int, content string, trackClear bool) {
	startX := (screenWidth - len(content)) / 2
	for i := 0; i < len(content); i++ {
		print(startX+i, startY, 1, 1, tcell.StyleDefault, rune(content[i]))
		if trackClear {
			coordinatesToClear = append(coordinatesToClear, &Coordinate{startX + i, startY})
		}
	}
	Screen.Show()
}

func clearScreen() {
	for _, coordinate := range coordinatesToClear {
		print(coordinate.x, coordinate.y, 1, 1, tcell.StyleDefault, ' ')
	}
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
		print(xOrigin+i, yOrigin, borderThickness, borderThickness, tcell.StyleDefault, upperBorder)
		print(xOrigin+i, yOrigin+height-1, borderThickness, borderThickness, tcell.StyleDefault, lowerBorder)
		// lower boundry
	}

	// side boundry
	for i := 1; i < height-1; i++ {
		print(xOrigin, yOrigin+i, borderThickness, borderThickness, tcell.StyleDefault, verticalBorder)
		print(xOrigin+width-1, yOrigin+i, borderThickness, borderThickness, tcell.StyleDefault, verticalBorder)
	}
}
