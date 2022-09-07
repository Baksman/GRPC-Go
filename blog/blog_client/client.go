package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/baksman/grpc/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NewBlog struct {
	Id       string
	AuthorId string
	Title    string
	Content  string
}

func main() {
	fmt.Println("hello  from average client")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("eror occured while dialing: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	blog := &blogpb.Blog{
		AuthorId: "Stephanie",
		Title:    "My first title",
		Content:  "Content of the first blog post",
	}

	createBlogResponse, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})

	if err != nil {
		log.Fatalf("error creating blog: %v", err)
	}
	fmt.Printf("blog has been created successfully %v\n", createBlogResponse)

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "63182400db29f41a8180a0e1",
	})

	if err != nil {
		log.Fatalf("error reading blog: %v\n", err)
	}

	fmt.Printf("we received blog id %v and title %v", res.Blog.GetAuthorId(), res.Blog.GetTitle())

	allBlogResponse, err := c.AllBlogs(context.Background(), &blogpb.BlogListRequest{})

	if err != nil {
		log.Fatalf("error updating blog: %v \n", err)
	}

	// allBlogs := []NewBlog{}

	// for _, blog := range allBlogResponse.Blogs {
	// 		blogItem = NewBlog
	// 	allBlogs = append(allBlogs, NewBlog)
	// }

	marshalRes, err := json.Marshal(allBlogResponse.Blogs)

fmt.Printf("response of all blogs: %v\n", string(marshalRes))

	// updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
	// 	Blog: &blogpb.Blog{
	// 		Id:       "63182024104883f205384216",
	// 		Title:    "This is the updated blog Title",
	// 		Content:  "ths content is fucking updated by me while creating grpc",
	// 		AuthorId: "Ibrahim Shehu",
	// 	},
	// })

	// if err != nil {
	// 	log.Fatalf("error updating blog: %v \n", err)
	// }
	// fmt.Printf("updated blog %v with success %v\n", updateRes.Blog.Title, updateRes.Blog.Content)
	// const deleteBlogId = "63182024104883f205384216"

	// deleteBlogRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
	// 	BlogId: deleteBlogId,
	// },
	// )

	// if err != nil {
	// 	log.Fatalf("error updating blog: %v \n", err)
	// }
	// fmt.Printf("updated blog %v with success %v\n", deleteBlogId, deleteBlogRes.Success)
}
