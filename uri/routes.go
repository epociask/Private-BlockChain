package uri

import (
	"./handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route



var routes = Routes{

	Route{
		"start",
		"GET",
		"/start",
		handlers.Start,
	},


	Route{
	//	//Response: If you have the block, return the JSON string of the specific block;
	//	//if you don't have the block, return HTTP 204: StatusNoContent;
	//	//if there's an error, return HTTP 500: InternalServerError.
	//	//Description: Return JSON string of a specific block to the downloader.
		"Get Block",
		"GET",
		"/block/{height}/{hash}",
		handlers.GetBlock,
	},

	Route{
		"Show",
		"GET",
		"/show",
		handlers.Show,

	},

	Route{
		//Description: Receive a block.
		"Receive Block",
		"POST",
		"/block/receive",
		handlers.ReceiveBlock,
	},

	Route{
	"Peer",
	"POST",
	"/peers",
	handlers.Register,


	},
	Route{
		"ShowPeers",
		"GET",
		"/showpeers",
		handlers.ShowPeers,
	},


}
