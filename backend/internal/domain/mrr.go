package domain

type MRR struct {
	New          float32
	Old          float32
	Reactivation float32
	Expansion    float32
	Contraction  float32
	Churn        float32
}

type TotalMRR struct {
	New          []float32
	Old          []float32
	Reactivation []float32
	Expansion    []float32
	Contraction  []float32
	Churn        []float32
	Total        []float32
}
