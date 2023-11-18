package main

import (
	"github.com/Muhammed-Rajab/go-blog/pkg/db"
	"github.com/Muhammed-Rajab/go-blog/pkg/routers"
	"github.com/Muhammed-Rajab/go-blog/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {

	db.Init("mongodb://localhost:27017")
	db.GetMDB().Connect()

	engine := gin.Default()
	engine.SetFuncMap(utils.GetTemplateFuncsMap())
	engine.LoadHTMLGlob("templates/*.html")

	root := engine.Group("/")
	routers.BlogRouter(root)

	engine.Run(":8000")

	// 	blogs := models.NewBlogs(db.GetMDB().BlogsCollection())
	// 	post := models.BlogModel{
	// 		Title: "Love and Hate: A Dichotomy in Our World",
	// 		Content: `![Love and Hate](https://example.com/love_hate_image.jpg)

	// In a world that spins on the axis of emotions, love and hate stand out as two powerful and opposing forces that shape the human experience. From the warmth of a compassionate embrace to the chilling winds of animosity, these emotions weave a complex tapestry that defines our relationships, societies, and, ultimately, our existence.

	// ## The Dance of Love

	// Love, often heralded as the most profound and positive emotion, has the capacity to bridge gaps, heal wounds, and create a sense of unity among diverse individuals. It manifests in myriad forms - the tender affection between parent and child, the bonds of friendship, and the romantic entanglements that ignite the human spirit.

	// ![Love](https://example.com/love_image.jpg)

	// ### Love's Impact on Society

	// The expression of love extends beyond the personal sphere, influencing the dynamics of entire societies. Acts of kindness, empathy, and compassion have the power to transform communities, fostering a sense of belonging and cooperation. Love, when harnessed collectively, can be a catalyst for positive change, breaking down barriers and promoting understanding.

	// ## The Shadows of Hate

	// Yet, in the tapestry of human emotions, hate casts a dark shadow. It is an emotion fueled by fear, prejudice, and misunderstanding, often leading to division, conflict, and suffering. Hate can manifest on personal, societal, and even global scales, leaving scars that endure through generations.

	// ![Hate](https://example.com/hate_image.jpg)

	// ### The Cycle of Hate

	// Hate begets hate in a vicious cycle that can seem difficult to break. Prejudices, stereotypes, and discrimination flourish in an environment poisoned by hatred. Recognizing and understanding the roots of hate is crucial to dismantling its destructive influence on individuals and societies.

	// ## Navigating the Spectrum

	// As we navigate the complex spectrum of human emotions, it becomes apparent that love and hate are not mutually exclusive. Individuals and societies grapple with these emotions in varying degrees, often experiencing both at different points in their journeys.

	// ### Striving for Balance

	// Understanding the coexistence of love and hate challenges us to strive for balance. By fostering love and empathy, we can build bridges that span the divides created by hatred. Education, open dialogue, and promoting diversity are essential tools in cultivating a world where love prevails over hate.

	// ## Conclusion

	// In the grand tapestry of life, love and hate are the vibrant threads that create a mosaic of human experience. Our challenge lies in recognizing the destructive potential of hate while harnessing the transformative power of love. As we weave our stories, let us choose the colors of compassion, understanding, and unity, working towards a world where love triumphs over hate.`,
	// 		Tags: []string{"love", "hate"},
	// 	}

	// if err := blogs.AddBlog(post); err != nil {
	// 	log.Fatal(err)
	// }
}
