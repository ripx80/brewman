package logger

// example for multiout logger
// f, _ := os.OpenFile("/tmp/brewman.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// mw := io.MultiWriter(os.Stderr, f)
// log.SetOutput(mw)


	//log.SetOutput(os.Stderr)

	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})
	contextLogger.Info("I'll be logged with common and other field")