package tests

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(suite.Suite))
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
