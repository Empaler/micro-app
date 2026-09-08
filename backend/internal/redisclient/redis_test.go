package redisclient

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
)

func TestClient_IncrementLookupsAndGetTop(t *testing.T) {
	server, err := miniredis.Run()
	require.NoError(t, err)

	client, err := NewClient(server.Addr(), 0)
	require.NoError(t, err)

	require.NoError(t, client.IncrementLookups(context.Background(), "movies", 10))
	require.NoError(t, client.IncrementLookups(context.Background(), "movies", 10))
	require.NoError(t, client.IncrementLookups(context.Background(), "movies", 20))
	require.NoError(t, client.IncrementLookups(context.Background(), "books", 30))

	topMovies, err := client.GetTopLookedUp(context.Background(), "movies", 5)
	require.NoError(t, err)
	require.Len(t, topMovies, 2)
	require.Equal(t, int64(10), topMovies[0].ID)
	require.EqualValues(t, 2, topMovies[0].Score)
	require.Equal(t, int64(20), topMovies[1].ID)
	require.EqualValues(t, 1, topMovies[1].Score)

	topBooks, err := client.GetTopLookedUp(context.Background(), "books", 5)
	require.NoError(t, err)
	require.Len(t, topBooks, 1)
	require.Equal(t, int64(30), topBooks[0].ID)
	require.EqualValues(t, 1, topBooks[0].Score)
}
