package main

import (
	// "context"

	"context"
	"fmt"
	"github.com/baksman/grpc/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type server struct {
	// calculatorpb.UnimplementedCalculatorServiceServer
}

var collection *mongo.Collection

type blogItem struct {
	ID       primitive.ObjectID `bson:"id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	// panic("error creating blog")
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(context.Background(), data)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "methods: %v", err)
		// log.Fatalf("error inserting blog: %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Errorf(codes.Internal, "cannot convert to oid: %v", err)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil

}

func (s *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	data := &blogItem{}
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "not a hexadecimal : %v", err)
	}
	filter := bson.D{primitive.E{Key: "_id", Value: oid}}

	mResult := collection.FindOneAndUpdate(context.Background(), filter, data)
	_ = mResult

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.String(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil

}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogId := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "not a hexadecimal : %v", err)
	}

	// create an empty
	data := &blogItem{}

	filter := bson.D{primitive.E{Key: "_id", Value: oid}}

	res := collection.FindOne(context.Background(), filter)

	err = res.Decode(data)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not found : %v", err)
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			AuthorId: string(data.AuthorID),
			Title:    data.Title,
			Content:  data.Content,
			Id:       data.ID.String(),
		},
	}, nil
}

func (s *server) AllBlogs(ctx context.Context, req *blogpb.BlogListRequest) (*blogpb.BlogListResponse, error) {
	dataList := []*blogItem{}
	blogs := []*blogpb.Blog{}

	filter := bson.D{}
	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not found : %v \n", err)
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		data := &blogItem{}

		cursor.Decode(data)
		dataList = append(dataList, data)
	}

	// err = res.Decode(&data)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, " error while decoding %v \n", err)
	// }

	fmt.Printf("lenght of %v", len(dataList))
	for _, item := range dataList {

		blogs = append(blogs, &blogpb.Blog{
			Id:       item.ID.String(),
			AuthorId: item.AuthorID,
			Title:    item.Title,
		})
	}

	return &blogpb.BlogListResponse{
		Blogs: blogs,
	}, nil

}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogId := req.GetBlogId()
	filter := bson.D{primitive.E{Key: "_id", Value: blogId}}

	deleteResult := collection.FindOneAndDelete(context.Background(), filter)

	if deleteResult.Err() != nil {
		return nil, status.Errorf(codes.NotFound, "blog with id %s not found", blogId)
	}
	return &blogpb.DeleteBlogResponse{
		Success: true,
	}, nil
	// return nil, status.Errorf(codes.Unimplemented, "method DeleteBlog not implemented")
}

func main() {
	// if we crash the go code we get file name and number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("blog server started")
	// connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("error connecting to database %v", err)
	}
	// create   database if not existing
	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to lsiten %v", err)
	}

	s := grpc.NewServer()

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("starting the server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()

	ch := make(chan os.Signal)

	signal.Notify(ch, os.Interrupt)

	<-ch

	fmt.Println("stoping the server")
	s.Stop()
	fmt.Println("closing mongodb server...")
	client.Disconnect(ctx)
	lis.Close()

	fmt.Println("server stopped")

}
