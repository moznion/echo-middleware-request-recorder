package echo_middleware_request_recorder

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/gob"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

type RequestRecorder struct {
	file *os.File
}

func NewRequestRecorder(outputFilePath string) (*RequestRecorder, error) {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return nil, err
	}

	return &RequestRecorder{
		file: file,
	}, nil
}

type Request struct {
	Method           string
	URL              *url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	Body             []byte
	ContentLength    int64
	TransferEncoding []string
	Host             string
	Form             url.Values
	PostForm         url.Values
	MultipartForm    *multipart.Form
	Trailer          http.Header
	RemoteAddr       string
	TLS              *tls.ConnectionState
}

func (s *RequestRecorder) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		body, _ := io.ReadAll(req.Body)
		_ = req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(&Request{
			Method:           req.Method,
			URL:              req.URL,
			Proto:            req.Proto,
			ProtoMajor:       req.ProtoMajor,
			ProtoMinor:       req.ProtoMinor,
			Header:           req.Header,
			Body:             body,
			ContentLength:    req.ContentLength,
			TransferEncoding: req.TransferEncoding,
			Host:             req.Host,
			Form:             req.Form,
			PostForm:         req.PostForm,
			MultipartForm:    req.MultipartForm,
			Trailer:          req.Trailer,
			RemoteAddr:       req.RemoteAddr,
			TLS:              req.TLS,
		})
		if err == nil {
			_, err := s.file.WriteString(base64.StdEncoding.EncodeToString(buf.Bytes()) + "\n")
			if err != nil {
				c.Logger().Errorf("request-recorder-middleware: failed to write a request log; %s", err)
			}
		} else {
			c.Logger().Errorf("request-recorder-middleware: failed to encode a request object to gob; %s", err)
		}

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
