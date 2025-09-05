package pulse

type Service struct {
	client    *Client
	analyzer  *Analyzer
	formatter *Formatter
}

func NewService() (*Service, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	return &Service{
		client:    client,
		analyzer:  NewAnalyzer(),
		formatter: NewFormatter(),
	}, nil
}

func (s *Service) GetClusterPulse(timeWindowMinutes int, podAmount int) (string, error) {
	pods, err := s.client.GetPodStatuses()
	if err != nil {
		return "", err
	}

	health := s.analyzer.AnalyzeClusterHealth(pods, timeWindowMinutes, podAmount)

	return s.formatter.FormatClusterHealth(health), nil
}
