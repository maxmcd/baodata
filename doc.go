// Baodata is a slow and simple datastore for gitbao
//
// The api uses REST-style endpoints to store data.
//
// Get:
//
//     resp, err := baodata.Get("/users/1")
//     // resp => []{"id":"1", "email":"git@gitbao.com"}
//
//     resp, err := baodata.Get("/users")
//     // resp => []{"id":"1", "email":"git@gitbao.com"}
//
// Put:
//
//     resp, err := baodata.Put("/users", baodata.Data{"email", "git@gitbao.com"})
//     // resp => []{"id":"2", "email":"git@gitbao.com"}
//
//     resp, err := baodata.Put("/users/2", baodata.Data{"email", "bao@gitbao.com"})
//     // resp => []{"id":"2", "email":"bao@gitbao.com"}
//
// Delete:
//
//     err := baodata.Delete("/users/1")
//
package baodata
