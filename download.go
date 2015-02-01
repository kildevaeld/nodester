package main

import (
	"errors"
	"io"
	"net/http"
)

type DownloadProgress struct {
	Progress int64
	Total    int64
}

func Download(url string, writer io.Writer) (chan int64, chan DownloadProgress, chan error) {
	done := make(chan int64)
	progress := make(chan DownloadProgress)
	errchan := make(chan error)

	go func() {
		res, err := http.Get(url)

		if err != nil {
			errchan <- err
			return
		}

		if res.StatusCode != 200 {
			// Go non 200
			errchan <- errors.New(res.Status)
			return
		}

		written, er := copy(writer, res.Body, func(p int64) {
			progress <- DownloadProgress{
				Progress: p,
				Total:    res.ContentLength,
			}
		})

		if er != nil {
			errchan <- er
			return
		}

		done <- written

	}()

	return done, progress, errchan
}

func DownloadSync(url string, writer io.Writer, progress func(DownloadProgress)) (written int64, err error) {

	done, prog, errchan := Download(url, writer)

loop:
	for {
		select {
		case n, ok := <-done:
			if !ok {
				err = errors.New("Undefined error")
				break loop
			}
			written = n
			break loop
		case p, ok := <-prog:
			if ok && progress != nil {
				progress(p)
			}

		case e, _ := <-errchan:
			err = e
			break loop
		}
	}
	return written, err

}

func copy(dest io.Writer, src io.Reader, progress func(p int64)) (written int64, err error) {

	buf := make([]byte, 32*1024)

	for {
		nr, er := src.Read(buf)

		if nr > 0 {
			nw, ew := dest.Write(buf[0:nr])

			if nw > 0 {
				written += int64(nw)
				progress(written)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}

		if er != nil {
			err = er
			break
		}

	}

	return written, err
}
