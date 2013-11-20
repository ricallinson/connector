package stackr

import (
	"net/http"
	"os"
)

/*
   Options for the static middleware. _Note: future options commented out._
*/
type staticOpt struct {
	Root string
	// MaxAge int64
	// Hidden bool
	// Redirect bool
	// Index bool
}

/*
   Static:

   Static file server with the given `root` path.

   Examples:

       stackr.CreateServer().Use(stackr.Static())
       stackr.CreateServer().Use(stackr.Static(map[string]string{"root": "./public"}))

   Options (not implemented yet):

       * `maxage` Browser cache maxAge in milliseconds. defaults to 0
       * `hidden` Allow transfer of hidden files. defaults to false
       * `redirect` Redirect to trailing "/" when the pathname is a dir. defaults to true
       * `index` Default file name, defaults to 'index.html'
*/
func Static(o ...map[string]string) func(req *Request, res *Response, next func()) {

	/*
	   If we got options use them.
	*/

	opt := staticOpt{}

	if len(o) == 1 {
		val := o[0]
		opt.Root = val["root"]
		// opt.MaxAge, _ = strconv.Atoi(val["maxage"])
		// opt.Hidden = val["hidden"] == "true"
		// opt.Redirect = val["redirect"] == "true"
		// opt.Index = val["index"] == "true"
	}

	/*
	   File Stat Cache.
	*/

	statCache := make(map[string]int)

	/*
	   The default loction of static files.
	*/

	root := "./public/"

	/*
	   If we were given a root use it.
	*/

	if len(opt.Root) > 0 {

		root = opt.Root

		/*
		   Add trailing slash if one is not there.
		*/

		if size := len(root); size > 1 && root[size-1] != '/' {
			root += "/"
		}
	}

	/*
	   Create a http.FileServer to server the files.
	*/

	fileServer := http.FileServer(http.Dir(root))

	/*
	   Return the Handle function.
	*/

	return func(req *Request, res *Response, next func()) {

		/*
		   Because http.FileServer serves directories and it's own 404 we
		   want to see if the file is really there before we hand of to it.
		   To do that we see if the file exists. If it doesn't, then we return quickly.

		   Question: Is this not really expensive?
		   Answer: It's not ideal. Writing a custom static server is on the todo list.
		*/

		filepath := root + req.Url

		/*
		   Check the stat cache as it's quicker than doing a stat on a file.
		*/

		if statCache[filepath] == -1 {
			next()
			return
		}

		/*
		   If the value of stat cache is 0 it means this is the first request for the filename.
		*/

		if statCache[filepath] == 0 {

			/*
			   Stat the filename.
			*/

			if stat, err := os.Stat(filepath); err != nil || stat.IsDir() == true {

				/*
				   If there is no file set stat cache to -1 and return.
				*/

				statCache[filepath] = -1

				next()
				return
			}

			/*
			   If there was a file set stat cache to 1 and let the FileServer serve it.
			*/

			statCache[filepath] = 1
		}

		/*
		   If we have to serve a file strip the matched Url and call ServeHTTP() on the fileServer.
		*/

		http.StripPrefix(req.MatchedUrl, fileServer).ServeHTTP(res.Writer, req.Request)

		/*
		   Now call End() to make sure we don't process any more middleware.
		*/

		res.End("")
	}
}
