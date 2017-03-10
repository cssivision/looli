// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package looli is a minimalist web framework for go
package main

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


Context supply some syntactic sugar

	    router := looli.Default()

	    router.Get("/query", func(c *looli.Context) {
	        id := c.Query("id")
	        name := c.DefaultQuery("name", "cssivision")
	        c.Status(200)
	        c.String("hello %s, %s\n", id, name)
	    })

	    router.Post("/form", func(c *looli.Context) {
	        name := c.DefaultPostForm("name", "somebody")
	        age := c.PostForm("age")
	        c.Status(200)
	        c.JSON(looli.JSON{
	            "name": name,
	            "age": age,
	        })
	    })

middleware useage

		func Logger() looli.HandlerFunc {
		    return func(c *looli.Context) {
		        t := time.Now()
		        // before request
		        c.Next()
		        // after request
		        latency := time.Since(t)
		        log.Print(latency)
		    }
		}

	    router := looli.New()

	    // global middleware
	    router.Use(Logger())
	    router.Get("/a", func(c *looli.Context) {
	        c.Status(200)
	        c.String("hello world!\n")
	    })

	    http.ListenAndServe(":8080", router)
*/
package looli
