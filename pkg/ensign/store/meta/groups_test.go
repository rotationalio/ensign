package meta_test

func (s *metaTestSuite) TestListGroups() {
	require := s.Require()
	require.False(s.store.ReadOnly())

}
