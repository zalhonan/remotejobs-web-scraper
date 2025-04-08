package service

import (
	"go.uber.org/zap"
)

func (s *service) CollectJobs() error {
	op := "internal.service.CollectJobs"

	for _, parser := range s.parsers {
		jobs, err := parser.ParseJobs()

		if err != nil {
			s.logger.Warn(
				op,
				zap.String("Parser retuned error while parsing jobs", parser.Name()),
				zap.Error(err),
			)
			continue
		}

		if len(jobs) == 0 {
			s.logger.Warn(
				op,
				zap.String("No jobs found while parsing", parser.Name()),
			)
			continue
		}

		saved, err := s.repository.SaveJobs(s.context, jobs)
		if err != nil {
			s.logger.Warn(
				op,
				zap.String("Error saving jobs from parser", parser.Name()),
				zap.Error(err),
			)
			continue
		}

		s.logger.Info(
			op,
			zap.String("Parser", parser.Name()),
			zap.Int("Jobs saved", saved),
		)
	}

	return nil
}
