package pkg

import (
	"time"
)

// Embedded is a base struct for all embeddedable structs.
type Embedded struct {
	Title           string `yaml:"title"`
	Slug            string `yaml:"slug"`
	Description     string `yaml:"description"`
	Content         string
	BannerPath      string `yaml:"banner_path"`
	RawContent      string
	Icon            string    `yaml:"icon"`
	CreatedAt       time.Time `yaml:"created_at"`
	UpdatedAt       time.Time `yaml:"updated_at"`
	X               float64
	Y               float64
	Z               float64
	TagSlugs        []string      `yaml:"tags"`
	PostSlugs       []string      `yaml:"posts"`
	ProjectSlugs    []string      `yaml:"projects"`
	EmploymentSlugs []string      `yaml:"employments"`
	Posts           []*Post       `yaml:"-" structgen:"PostSlugs"`
	Tags            []*Tag        `yaml:"-" structgen:"TagSlugs"`
	Projects        []*Project    `yaml:"-" structgen:"ProjectSlugs"`
	Employments     []*Employment `yaml:"-" structgen:"EmploymentSlugs"`
}

// Post is a post with all its projects and tags.
type Post struct {
	Embedded
}

// Project is a project with all its posts and tags.
type Project struct {
	Embedded
}

// Tag is a tag with all its posts and projects.
type Tag struct {
	Embedded
}

// Employment is an employment of a tag.
type Employment struct {
	Embedded
}

// Posts is a collection of posts.
var Posts = []*Post{
	{
		Embedded: Embedded{
			Title:        "Getting Started with Go",
			Slug:         "getting-started-with-go",
			Description:  "A beginner's guide to Go programming",
			Content:      "Go is a statically typed compiled language designed at Google...",
			BannerPath:   "/images/go-banner.jpg",
			CreatedAt:    time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:    time.Date(2023, 5, 16, 9, 30, 0, 0, time.UTC),
			TagSlugs:     []string{"go", "programming", "beginners"},
			ProjectSlugs: []string{"go-tutorial"},
		},
	},
	{
		Embedded: Embedded{
			Title:        "Advanced Go Concurrency",
			Slug:         "advanced-go-concurrency",
			Description:  "Deep dive into Go's concurrency model",
			Content:      "Go's goroutines and channels offer a powerful concurrency model...",
			BannerPath:   "/images/concurrency-banner.jpg",
			CreatedAt:    time.Date(2023, 6, 10, 14, 0, 0, 0, time.UTC),
			UpdatedAt:    time.Date(2023, 6, 12, 11, 45, 0, 0, time.UTC),
			TagSlugs:     []string{"go", "concurrency", "advanced"},
			ProjectSlugs: []string{"go-tutorial"},
		},
	},
}

// Tags is a list of tags.
var Tags = []*Tag{
	{
		Embedded: Embedded{
			Title:       "Go",
			Slug:        "go",
			Description: "The Go programming language",
			Icon:        "go-icon",
			CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			PostSlugs:   []string{"getting-started-with-go", "advanced-go-concurrency"},
		},
	},
	{
		Embedded: Embedded{
			Title:       "Programming",
			Slug:        "programming",
			Description: "General programming topics",
			Icon:        "code-icon",
			CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			PostSlugs:   []string{"getting-started-with-go"},
		},
	},
	{
		Embedded: Embedded{
			Title:       "Beginners",
			Slug:        "beginners",
			Description: "Content for beginners",
			Icon:        "beginner-icon",
			CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			PostSlugs:   []string{"getting-started-with-go"},
		},
	},
	{
		Embedded: Embedded{
			Title:       "Concurrency",
			Slug:        "concurrency",
			Description: "Parallel and concurrent programming topics",
			Icon:        "concurrency-icon",
			CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			PostSlugs:   []string{"advanced-go-concurrency"},
		},
	},
	{
		Embedded: Embedded{
			Title:       "Advanced",
			Slug:        "advanced",
			Description: "Advanced programming topics",
			Icon:        "advanced-icon",
			CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			PostSlugs:   []string{"advanced-go-concurrency"},
		},
	},
}

// Projects is a collection of projects.
var Projects = []*Project{
	{
		Embedded: Embedded{
			Title:       "Go Tutorial",
			Slug:        "go-tutorial",
			Description: "Comprehensive Go tutorial from basics to advanced",
			BannerPath:  "/images/tutorial-banner.jpg",
			CreatedAt:   time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC),
			PostSlugs:   []string{"getting-started-with-go", "advanced-go-concurrency"},
			TagSlugs:    []string{"go", "programming", "beginners", "advanced"},
		},
	},
}

var Employments = []*Employment{
	{
		Embedded: Embedded{
			Title:       "Google",
			Slug:        "google",
			Description: "Software Engineer at Google",
			Icon:        "google-icon",
			CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			TagSlugs:    []string{"go", "programming"},
		},
	},
	{
		Embedded: Embedded{
			Title:       "Microsoft",
			Slug:        "microsoft",
			Description: "Software Engineer at Microsoft",
			Icon:        "microsoft-icon",
			CreatedAt:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			TagSlugs:    []string{"programming"},
		},
	},
}
