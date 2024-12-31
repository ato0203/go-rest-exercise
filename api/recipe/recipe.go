package recipe

import (
	"bukeuw/recipe/pkg/recipes"
	"net/http"
  "github.com/gosimple/slug"

	"github.com/gin-gonic/gin"
)

type recipeStore interface {
  Add(name string, recipe recipes.Recipe) error
  List() (map[string]recipes.Recipe, error)
  Get(name string) (recipes.Recipe, error)
  Update(name string, recipe recipes.Recipe) error
  Remove(name string) error
}

type RecipeHandler struct {
  store recipeStore
}

func NewRecipeHandler(s recipeStore) *RecipeHandler {
  return &RecipeHandler{
    store: s,
  }
}

func (h RecipeHandler) ListRecipe(ctx *gin.Context) {
  recipes, err := h.store.List()
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })
  }
  ctx.JSON(http.StatusOK, gin.H{
    "recipes": recipes,
  })
}

func (h RecipeHandler) CreateRecipe(ctx *gin.Context) {
  var recipe recipes.Recipe
  if err := ctx.ShouldBindJSON(&recipe); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "error": err.Error(),
    })
    return
  }
  
  id := slug.Make(recipe.Name)
  h.store.Add(id, recipe)
  ctx.JSON(http.StatusCreated, gin.H{
    "recipe": recipe,
  })
}

func (h RecipeHandler) GetRecipe(ctx *gin.Context) {
  recipeId := ctx.Param("id")
  recipe, err := h.store.Get(recipeId)
  if err != nil {
    ctx.JSON(http.StatusNotFound, gin.H{
      "error": err.Error(),
    })
    return
  }
  ctx.JSON(http.StatusOK, gin.H{
    "recipe": recipe,
  })
}

func (h RecipeHandler) UpdateRecipe(ctx *gin.Context) {
  var recipe recipes.Recipe
  if err := ctx.ShouldBindJSON(&recipe); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "error": err.Error(),
    })
    return
  }

  recipeId := ctx.Param("id")
  err := h.store.Update(recipeId, recipe)
  if err != nil {
    if err == recipes.NotFoundErr {
      ctx.JSON(http.StatusNotFound, gin.H{
        "error": err.Error(),
      })
      return
    }
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })
    return
  }
  ctx.JSON(http.StatusOK, gin.H{
    "recipe": recipe,
  })
}

func (h RecipeHandler) DeleteRecipe(ctx *gin.Context) {
  recipeId := ctx.Param("id")
  err := h.store.Remove(recipeId)
  if err != nil {
    if err == recipes.NotFoundErr {
      ctx.JSON(http.StatusNotFound, gin.H{
        "error": err.Error(),
      })
      return
    }
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })
    return
  }
  ctx.JSON(http.StatusOK, gin.H{
    "message": "recipe deleted",
    "recipe_id": recipeId,
  })
}
