package books

import (
	"net/http"
	"strconv"

	"movie-api/internal/redisclient"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Router struct {
	service *BookService
	redis   *redisclient.Client
}

func RegisterRoutes(g *gin.Engine, db *sqlx.DB, redisClient *redisclient.Client) {
	adapter := NewPostgresAdapter(db)
	service := NewBookService(adapter)
	router := &Router{service: service, redis: redisClient}
	router.registerRoutes(g)
}

func (r *Router) registerRoutes(g *gin.Engine) {
	api := g.Group("/api")
	{
		api.GET("/books", r.listBooks)
		api.GET("/books/:id", r.getBook)
		api.GET("/books/most-looked-up", r.mostLookedUpBooks)
		api.POST("/books", r.createBook)
		api.PUT("/books/:id", r.updateBook)
		api.DELETE("/books/:id", r.deleteBook)
	}
}

// @Summary List all books
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/books [get]
func (r *Router) listBooks(c *gin.Context) {
	books, err := r.service.GetAllBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": books})
}

// @Summary Get a book by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/books/{id} [get]
func (r *Router) getBook(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return
	}

	book, err := r.service.GetBook(c.Request.Context(), id)
	if err != nil {
		if err == ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if r.redis != nil {
		_ = r.redis.IncrementLookups(c.Request.Context(), "books", id)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": book})
}

// @Summary Get the most looked up books
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/books/most-looked-up [get]
func (r *Router) mostLookedUpBooks(c *gin.Context) {
	if r.redis == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []redisclient.PopularItem{}})
		return
	}

	items, err := r.redis.GetTopLookedUp(c.Request.Context(), "books", 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// @Summary Create a new book
// @Tags books
// @Accept json
// @Produce json
// @Param book body Book true "Book data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/books [post]
func (r *Router) createBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := r.service.CreateBook(c.Request.Context(), &book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": book})
}

// @Summary Update a book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body Book true "Book data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/books/{id} [put]
func (r *Router) updateBook(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return
	}

	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	book.ID = id

	if err := r.service.UpdateBook(c.Request.Context(), &book); err != nil {
		if err == ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "book not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": book})
}

// @Summary Delete a book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Router /api/books/{id} [delete]
func (r *Router) deleteBook(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return
	}

	if err := r.service.DeleteBook(c.Request.Context(), id); err != nil {
		if err == ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
