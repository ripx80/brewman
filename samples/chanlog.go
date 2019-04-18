
		logc := make(chan string)
		done := make(chan error, 2)
		defer close(done)

		go func() {
			done <- masher.Mash(logc)
		}()

		go func() {
			done <- func() error {
				for {
					j, more := <-logc
					if more {
						log.Info(j)
					} else {
						return nil
					}
				}
			}()
		}()

		var stopped bool
		for i := 0; i < cap(done); i++ {
			if err := <-done; err != nil {
				log.Error(err)
			}
			if !stopped {
				stopped = true
				close(logc)
				log.Info("Close")
			}
		}
		log.Info("Q")
		//https://github.com/heptio/workgroup