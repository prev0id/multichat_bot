// Package desc provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package desc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Leave the twitch chat for the specified channel
	// (DELETE /twitch/{chat})
	LeaveTwitchChat(w http.ResponseWriter, r *http.Request, chat string)
	// Join the twitch chat for the specified channel
	// (POST /twitch/{chat})
	JoinTwitchChat(w http.ResponseWriter, r *http.Request, chat string)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Leave the twitch chat for the specified channel
// (DELETE /twitch/{chat})
func (_ Unimplemented) LeaveTwitchChat(w http.ResponseWriter, r *http.Request, chat string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Join the twitch chat for the specified channel
// (POST /twitch/{chat})
func (_ Unimplemented) JoinTwitchChat(w http.ResponseWriter, r *http.Request, chat string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// LeaveTwitchChat operation middleware
func (siw *ServerInterfaceWrapper) LeaveTwitchChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "chat" -------------
	var chat string

	err = runtime.BindStyledParameterWithLocation("simple", false, "chat", runtime.ParamLocationPath, chi.URLParam(r, "chat"), &chat)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "chat", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.LeaveTwitchChat(w, r, chat)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// JoinTwitchChat operation middleware
func (siw *ServerInterfaceWrapper) JoinTwitchChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "chat" -------------
	var chat string

	err = runtime.BindStyledParameterWithLocation("simple", false, "chat", runtime.ParamLocationPath, chi.URLParam(r, "chat"), &chat)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "chat", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.JoinTwitchChat(w, r, chat)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/twitch/{chat}", wrapper.LeaveTwitchChat)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/twitch/{chat}", wrapper.JoinTwitchChat)
	})

	return r
}

type LeaveTwitchChatRequestObject struct {
	Chat string `json:"chat"`
}

type LeaveTwitchChatResponseObject interface {
	VisitLeaveTwitchChatResponse(w http.ResponseWriter) error
}

type LeaveTwitchChat200Response struct {
}

func (response LeaveTwitchChat200Response) VisitLeaveTwitchChatResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type LeaveTwitchChat500Response struct {
}

func (response LeaveTwitchChat500Response) VisitLeaveTwitchChatResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

type JoinTwitchChatRequestObject struct {
	Chat string `json:"chat"`
}

type JoinTwitchChatResponseObject interface {
	VisitJoinTwitchChatResponse(w http.ResponseWriter) error
}

type JoinTwitchChat200Response struct {
}

func (response JoinTwitchChat200Response) VisitJoinTwitchChatResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type JoinTwitchChat500Response struct {
}

func (response JoinTwitchChat500Response) VisitJoinTwitchChatResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Leave the twitch chat for the specified channel
	// (DELETE /twitch/{chat})
	LeaveTwitchChat(ctx context.Context, request LeaveTwitchChatRequestObject) (LeaveTwitchChatResponseObject, error)
	// Join the twitch chat for the specified channel
	// (POST /twitch/{chat})
	JoinTwitchChat(ctx context.Context, request JoinTwitchChatRequestObject) (JoinTwitchChatResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHttpHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHttpMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// LeaveTwitchChat operation middleware
func (sh *strictHandler) LeaveTwitchChat(w http.ResponseWriter, r *http.Request, chat string) {
	var request LeaveTwitchChatRequestObject

	request.Chat = chat

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.LeaveTwitchChat(ctx, request.(LeaveTwitchChatRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "LeaveTwitchChat")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(LeaveTwitchChatResponseObject); ok {
		if err := validResponse.VisitLeaveTwitchChatResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// JoinTwitchChat operation middleware
func (sh *strictHandler) JoinTwitchChat(w http.ResponseWriter, r *http.Request, chat string) {
	var request JoinTwitchChatRequestObject

	request.Chat = chat

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.JoinTwitchChat(ctx, request.(JoinTwitchChatRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "JoinTwitchChat")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(JoinTwitchChatResponseObject); ok {
		if err := validResponse.VisitJoinTwitchChatResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
