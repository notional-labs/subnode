package rpc

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewRpcProxy takes target host and creates a reverse proxy
func NewRpcProxy(targetHost string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	log.Printf("forwarding to -> %s\n", targetHost)

	proxy := httputil.NewSingleHostReverseProxy(target)

	//originalDirector := proxy.Director
	//proxy.Director = func(req *http.Request) {
	//	originalDirector(req)
	//	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
	//}

	//proxy.ModifyResponse = modifyResponse()
	//proxy.ErrorHandler = errorHandler()

	return proxy, nil
}

//func errorHandler() func(http.ResponseWriter, *http.Request, error) {
//	return func(w http.ResponseWriter, req *http.Request, err error) {
//		fmt.Printf("Got error while modifying response: %v \n", err)
//		return
//	}
//}

//func modifyResponse() func(*http.Response) error {
//	return func(resp *http.Response) error {
//		return errors.New("response body is invalid")
//	}
//}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		proxy.ServeHTTP(w, r)
	}
}

func StartRpcServer() {
	// initialize a reverse proxy and pass the actual backend server url here
	rpcProxy, err := NewRpcProxy("https://rpc-osmosis-ia.cosmosia.notional.ventures:443")

	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(rpcProxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
	//log.Fatal(http.ListenAndServe(":8080", rpcProxy))
}
