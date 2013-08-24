package stackr

/*
    The middleware type.
*/

type middleware struct {
    Route string
    Handle func(*Request, *Response, func())
}