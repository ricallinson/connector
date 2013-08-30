package stackr

import(
    "net/http"
)

/*
    A HTTP Request.
*/
type Request struct {
    Raw *http.Request
    Url string
    MatchedUrl string
    OriginalUrl string
}

/*
    Returns a new Request.
*/

func createRequest(raw *http.Request) (*Request) {

    /*
        Create a new Request.
    */

    req := new(Request)

    /*
        Set the source http.Request so it can be accessed later.
    */

    req.Raw = raw

    /*
        Set the Url for easy access.

        Note: this value may be changed by the stack.handle() function.
    */

    req.Url = raw.RequestURI

    /*
        Set the Url for easy access.

        Note: this value should never change over the life time of the request.
    */

    req.OriginalUrl = raw.RequestURI

    /*
        Return the finished stack.Request.
    */

    return req
}