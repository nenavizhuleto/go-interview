package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Headers            http.Header
	UserAgent          string
	Data               string
	Method             string
	Insecure           bool
	Url                *url.URL
	ControlOutput      io.Writer
	ResponseBodyOutput io.Writer
}

func main() {

	if err := CreateCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
func Execute(c *Config) error {
	var r io.Reader
	var tlsConfig *tls.Config

	if c.Data != "" {
		r = bytes.NewBufferString(c.Data)
	}

	if c.Insecure {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	request, err := http.NewRequest(c.Method, c.Url.String(), r)
	if err != nil {
		return err
	}

	if c.UserAgent != "" {
		request.Header.Set("User-Agent", c.UserAgent)
	}

	for key, values := range c.Headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	requestBuilder := &wrappedBuilder{
		prefix: ">",
	}

	requestBuilder.Printf("%v %v", request.Method, request.URL.String())
	requestBuilder.WriteHeaders(request.Header)
	requestBuilder.Println()

	if _, err := io.Copy(c.ControlOutput, strings.NewReader(requestBuilder.String())); err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			fmt.Println("failed to close response body")
		}
	}()

	responseBuilder := &wrappedBuilder{
		prefix: "<",
	}

	responseBuilder.Printf("%v %v", response.Proto, response.Status)
	responseBuilder.WriteHeaders(response.Header)
	responseBuilder.Printf("")
	responseBuilder.Println()

	if _, err := io.Copy(c.ControlOutput, strings.NewReader(responseBuilder.String())); err != nil {
		return err
	}

	_, err = io.Copy(c.ResponseBodyOutput, response.Body)
	return err
}

type wrappedBuilder struct {
	prefix string
	strings.Builder
}

func (w *wrappedBuilder) WriteHeaders(headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			w.Printf("%v: %v", key, value)
		}
	}
}

func (w *wrappedBuilder) Println() {
	w.WriteString("\n")
}

func (w *wrappedBuilder) Printf(s string, a ...any) {
	w.WriteString(fmt.Sprintf("%v %v\n", w.prefix, fmt.Sprintf(s, a...)))
}

func CreateCommand() *cobra.Command {
	config := &Config{
		Headers:            map[string][]string{},
		ResponseBodyOutput: os.Stdout,
		ControlOutput:      os.Stdout,
	}

	headers := make([]string, 0, 255)

	command := &cobra.Command{
		Use:     `gurl URL`,
		Short:   `gurl is an HTTP client`,
		Long:    `gurl is an HTTP client for a tutorial of learing golang`,
		Args:    ArgsValidator(config),
		PreRunE: OptionsValidator(config, headers),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Execute(config)
		},
	}

	command.PersistentFlags().StringSliceVarP(&headers, "headers", "H", nil, `custom headers headers to be sent with the request, headers are separated by "," as in "HeaderName: Header content,OtherHeader: Some other value"`)
	command.PersistentFlags().StringVarP(&config.UserAgent, "user-agent", "u", "gurl", "the user agent to be used for requests")
	command.PersistentFlags().StringVarP(&config.Data, "data", "d", "", "data to be sent as the request body")
	command.PersistentFlags().StringVarP(&config.Method, "method", "m", http.MethodGet, "HTTP method to be used for the request")
	command.PersistentFlags().BoolVarP(&config.Insecure, "insecure", "k", false, "allows insecure server connections over HTTPS")

	return command
}

func OptionsValidator(c *Config, headers []string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, h := range headers {
			if name, value, found := strings.Cut(h, ":"); found {
				c.Headers.Add(strings.TrimSpace(name), strings.TrimSpace(value))
			} else {
				fmt.Printf("header is not a valid http header separated by `:`, value was: [%v]", h)
			}
		}

		return nil
	}
}

func ArgsValidator(c *Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if l := len(args); l != 1 {
			fmt.Printf("you must provide a single URL to be called but you provided %v\n", l)
			os.Exit(1)
		}

		u, err := url.Parse(args[0])
		if err != nil {
			fmt.Printf("the URL provided is invalid: %v", args[0])
		}

		c.Url = u

		return nil
	}
}
