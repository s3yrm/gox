package main

import (
	"flag"
	"log"
	"os"

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

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, logfile)
	Logger = zerolog.New(multi).With().Timestamp().Logger()
	Logger = Logger.With().Caller().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Remove this to pass filename from commandline //
	file := File{
		Name:   "Goxfile",
		Logger: Logger,
	}

	err = file.Open()
	if err != nil {
		Logger.Error().Err(err).Msg("Error working with file.")
	}

}
