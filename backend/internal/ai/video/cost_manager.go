package video

type CostManager struct {
	monthlyBudget float64
	usedBudget    float64
	userQuotas    map[string]int // userID -> remaining generations
}

func NewCostManager() *CostManager {
	return &CostManager{
		monthlyBudget: 0.0, // 0 = مجاني فقط
		userQuotas: map[string]int{
			"free_tier":    5, // 5 فيديوهات مجانية/شهر
			"basic_tier":   20,
			"premium_tier": 100,
		},
	}
}

func (c *CostManager) CanGenerateVideo(userID, tier string) bool {
	quota, exists := c.userQuotas[tier]
	if !exists {
		quota = c.userQuotas["free_tier"]
	}

	// تحقق من الحصة
	return quota > 0
}

func (c *CostManager) RecordGeneration(userID, tier string, cost float64) {
	if quota, exists := c.userQuotas[tier]; exists && quota > 0 {
		c.userQuotas[tier] = quota - 1
	}
	c.usedBudget += cost
}
