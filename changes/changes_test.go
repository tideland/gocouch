package changes_test

// TestChanges tests retrieving changes.
// func TestChanges(t *testing.T) {
// 	assert := audit.NewTestingAssertion(t, true)
// 	cdb, cleanup := prepareFilledDatabase("changes", assert)
// 	defer cleanup()

	// Simple changes access.
// 	crs := cdb.Changes()
// 	assert.True(crs.IsOK())
// 	assert.True(crs.ResultsLen() > 0)

// 	crs.ResultsDo(func(id, sequence string, deleted bool, revisions ...string) error {
// 		assert.Logf("%v: %v / %v / %v", id, sequence, deleted, revisions)
// 		return nil
// 	})
// }

