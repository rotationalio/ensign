package meta_test

func (s *metaTestSuite) TestListGroups() {
	require := s.Require()
	require.False(s.store.ReadOnly())

}

func (s *metaTestSuite) TestSameGroupName() {
	// Two different projects should be able to store a group with the same name without
	// conflict with or without group IDs specified by the user.
}
