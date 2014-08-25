tellmeabout
===========

Experimenting with Go's [reflect](http://golang.org/pkg/reflect/) package. TellMeAbout is a method that takes an object of any type and describes it recursively.

For example, given the structures:

```go
type Post struct {
	Id   int
	Name string
	User *User
}

type User struct {
	Id   int
	Name string
}
```

Doing the following:

```go
TellMeAbout(&Post{
	1, "About me", &User{1, "James"},
})
```

Will output:

```
You've passed a pointer to Post, with fields:

   - int
   - string
   - pointer to User, with fields:

      - int
      - string

```
