package service

import "fmt"

func (s *service) CollectJobs() error {
	for _, parser := range s.parsers {
		jobs, err := parser.ParseJobs()

		if err != nil {
			fmt.Printf("Parser %v retuned error while parsing jobs: %v\n", parser.Name(), err)
			continue
		}

		if len(jobs) == 0 {
			fmt.Printf("No jobs found while parsing %v\n", parser.Name())
			continue
		}

		saved, err := s.repository.SaveJobs(s.context, jobs)
		if err != nil {
			fmt.Printf("Error saving jobs from parser %v: %v\n", parser.Name(), err)
			continue
		}

		fmt.Printf("Collected and saved to pepository %d jobs from parser %v\n", saved, parser.Name())
	}

	return nil
}
