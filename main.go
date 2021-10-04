package main

import (
	"flag"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	logfile, err := os.Create("goxlog.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer logfile.Close()

	Logger = zerolog.New(logfile).With().Timestamp().Logger()
	Logger = Logger.With().Caller().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// begin setting up the UI
	scr, err := tcell.NewScreen()
	if err != nil {
		Logger.Fatal().Msg("Cannot Create screen")
	}

	err = scr.Init()
	if err != nil {
		Logger.Fatal().Msg("Cannot Initialize screen")
	}

	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	scr.SetStyle(style)
	scr.DisableMouse()

	scr.Clear()

	scrWidth, scrHeight := scr.Size()

	// Draw line begin indicators
	for row, col := 0, scrHeight; col >= 0; col-- {
		scr.SetContent(row, col, '~', nil, style)
	}

	// Place cursor at the top left and show screen
	scr.ShowCursor(0, 0)

	// row and col to keep track of cursor on screen; init at (0,0)
	row, col := 0, 0

	// quit the application
	quit := func() {
		scr.Fini()
		os.Exit(0)
	}

	for {
		// Show screen
		scr.Show()

		// Poll events
		ev := scr.PollEvent()

		// Process events
		switch ev := ev.(type) {
		case *tcell.EventResize:
			scrWidth, scrHeight = scr.Size()
			scr.Sync()
			for firstCol := scrHeight; firstCol >= 0; firstCol-- {
				scr.SetContent(0, firstCol, '~', nil, style)
			}
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
			if ev.Key() == tcell.KeyRune {
				// If the row rows more than the width of the screen, insert cursor and content in new line, beginning at 0.
				if row >= scrWidth {
					col++
					row = 0
					Logger.Debug().Msgf("Line break; row: %d, Col: %d", row, col)
				}
				scr.SetContent(row, col, ev.Rune(), nil, style)
				scr.ShowCursor(row+1, col) // cursor will always be one rune ahead.
				row++

			}
			if ev.Key() == tcell.KeyBackspace2 {
				scr.SetContent(row-1, col, rune(tcell.KeyNUL), nil, style) // rune appearing before the cursor is reset.
				row--

				if row <= 0 && col <= 0 { // If the cursor is at origin, do not do anything.
					Logger.Debug().Msg("Cursor at origin.")
					row, col = 0, 0
					scr.ShowCursor(0, 0)
					break
				} else if row <= 0 && col > 0 { // If the content is multiple lines long and cursor is at first row, place it 
				Logger.Debug().Msgf("Line cleared, moving to previous line; row: %d, col: %d", row, col)
					row = scrWidth // at the end of the previous line.
					col--
					scr.ShowCursor(row, col)
					break
				} else {
					scr.ShowCursor(row, col)
				}
			}
		}
	}
}
