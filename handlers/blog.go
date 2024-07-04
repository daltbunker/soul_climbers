package handlers

import (
	"fmt"
	"net/http"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
)

func GetBlogs(w http.ResponseWriter, r *http.Request) []types.Blog {
	blogs, err := db.GetAllBlogs(r)
	if err != nil {
		HandleServerError(w, err)
	}
	return blogs
}

func HandleNewBlog(w http.ResponseWriter, r *http.Request) {
	blog := types.Blog{}
	blog.CreatedBy = "jonnyX"
	blog.Title = "This is just another test article"
	blog.Body = `Lead researcher Mr. Thornbuckle proves the effectiveness of the participants across
	Computers aren't random. On the contrary, hardware designers work very hard to make sure computers run every program the same way every time. So when a program does
	need random numbers, that requires extra effort. Traditionally, computer scientists and programming languages have distinguished between two different kinds of random numbers: 
	<br><br>
	<h3 class="text-2xl">Traditional Computers</h3>
	statistical and cryptographic randomness. In Go, those are provided by math/rand and crypto/rand, respectively. This post is about how Go 1.22 brings the two closer together,
	by using a cryptographic random number source in math/rand (as well as math/rand/v2, as mentioned in our previous post). The result is better randomness and far less 
	<br><br>
	<h3 class="text-2xl">Cryptographic</h3>
	damage when developers accidentally use math/rand instead of crypto/rand.
	 `

	savedBlog, err := db.NewBlog(r, blog)
	if err != nil {
		HandleClientError(w, err)
	}

	// renderComponent(w, "home", "blog", savedBlog)
	fmt.Fprintf(w, "<h2>%s</h2> author %s <p>%s</p>", savedBlog.Title, savedBlog.CreatedBy, savedBlog.Body)
}
