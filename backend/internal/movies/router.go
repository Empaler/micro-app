package movies

import (
	"net/http"
	"strconv"

	"movie-api/internal/redisclient"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Router struct {
	service *MovieService
	redis   *redisclient.Client
}

func RegisterRoutes(g *gin.Engine, db *sqlx.DB, redisClient *redisclient.Client) {
	adapter := NewPostgresAdapter(db)
	service := NewMovieService(adapter)
	router := &Router{service: service, redis: redisClient}
	router.registerRoutes(g)
}

func (r *Router) registerRoutes(g *gin.Engine) {
	api := g.Group("/api")
	{
		api.GET("/movies", r.listMovies)
		api.GET("/movies/:id", r.getMovie)
		api.GET("/movies/most-looked-up", r.mostLookedUpMovies)
		api.POST("/movies", r.createMovie)
		api.PUT("/movies/:id", r.updateMovie)
		api.DELETE("/movies/:id", r.deleteMovie)
	}
}

// @Summary List all movies
// @Tags movies
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/movies [get]
func (r *Router) listMovies(c *gin.Context) {
	movies, err := r.service.GetAllMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": movies})
}

// @Summary Get a movie by ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/movies/{id} [get]
func (r *Router) getMovie(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return
	}

	movie, err := r.service.GetMovie(c.Request.Context(), id)
	if err != nil {
		if err == ErrMovieNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if r.redis != nil {
		_ = r.redis.IncrementLookups(c.Request.Context(), "movies", id)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movie})
}

// @Summary Get the most looked up movies
// @Tags movies
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/movies/most-looked-up [get]
func (r *Router) mostLookedUpMovies(c *gin.Context) {
	if r.redis == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []redisclient.PopularItem{}})
		return
	}

	items, err := r.redis.GetTopLookedUp(c.Request.Context(), "movies", 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// @Summary Create a new movie
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body Movie true "Movie data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/movies [post]
func (r *Router) createMovie(c *gin.Context) {
	var movie Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := r.service.CreateMovie(c.Request.Context(), &movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": movie})
}

// @Summary Update a movie
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body Movie true "Movie data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/movies/{id} [put]
func (r *Router) updateMovie(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return
	}

	var movie Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	movie.ID = id

	if err := r.service.UpdateMovie(c.Request.Context(), &movie); err != nil {
		if err == ErrMovieNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "movie not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movie})
}

// @Summary Delete a movie
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Router /api/movies/{id} [delete]
func (r *Router) deleteMovie(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return
	}

	if err := r.service.DeleteMovie(c.Request.Context(), id); err != nil {
		if err == ErrMovieNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
