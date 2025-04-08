package service

import (
	"go.uber.org/zap"
)

func (s *service) CollectJobs() error {
	for _, parser := range s.parsers {
		jobs, err := parser.ParseJobs()

		if err != nil {
			s.logger.Warn(
				"Parser retuned error while parsing jobs",
				zap.String("Parser", parser.Name()),
				zap.Error(err),
			)
			continue
		}

		if len(jobs) == 0 {
			s.logger.Warn(
				"No jobs found while parsing",
				zap.String("Parser", parser.Name()),
			)
			continue
		}

		saved, err := s.repository.SaveJobs(s.context, jobs)
		if err != nil {
			s.logger.Warn(
				"Error saving jobs from parser",
				zap.String("Parser", parser.Name()),
				zap.Error(err),
			)
			continue
		}

		s.logger.Info(
			"Parsing successfully completed",
			zap.String("Parser", parser.Name()),
			zap.Int("Jobs saved", saved),
		)
	}

	return nil
}
