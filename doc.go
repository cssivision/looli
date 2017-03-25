// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package looli is a minimalist web framework for go

router usage
    router := looli.Default()

    v1 := router.Prefix("/v1")
    v1.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world version1\n")
    })

    v2 := router.Prefix("/v2")
    v2.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world version2\n")
    })

    router.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world!\n")
    })
*/
package looli
