package storage

import (
	serviceErrors "github.com/sergeysynergy/metricser/internal/service/errors"
	"log"
)

func (s *Storage) snapShotRestore() (err error) {
	if s.fileRepo == nil {
		return serviceErrors.ErrFileStoreNotDefined
	}

	prm, err := s.fileRepo.JustReadMetrics()
	if err != nil {
		return err
	}

	err = s.repo.PutMetrics(prm)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Metrics has been restored from filestore")
	return nil
}
