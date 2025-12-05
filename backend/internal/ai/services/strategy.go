package services

type StrategyService struct {
    textProvider TextProvider
}

func NewStrategyService(provider TextProvider) *StrategyService {
    return &StrategyService{textProvider: provider}
}

// GenerateGrowthStrategy إنشاء استراتيجية نمو
func (s *StrategyService) GenerateGrowthStrategy(businessType, goals string) (string, error) {
    prompt := fmt.Sprintf(`
    Create a 90-day digital growth strategy for a %s business.
    
    Business Goals: %s
    
    The strategy should include:
    
    1. TARGET AUDIENCE:
    - Ideal customer profile
    - Persona development
    - Pain points and solutions
    
    2. CONTENT STRATEGY:
    - Content pillars
    - Editorial calendar
    - Content formats and channels
    
    3. SOCIAL MEDIA PLAN:
    - Platform selection
    - Posting frequency
    - Engagement tactics
    
    4. SEO STRATEGY:
    - Keyword research
    - On-page optimization
    - Link building
    
    5. PAID ADVERTISING:
    - Platform recommendations
    - Budget allocation
    - Targeting options
    
    6. METRICS & KPIs:
    - Key performance indicators
    - Measurement tools
    - Reporting schedule
    
    7. 90-DAY TIMELINE:
    - Month 1: Foundation
    - Month 2: Execution
    - Month 3: Optimization
    
    Make it specific, measurable, and actionable.
    `, businessType, goals)
    
    return s.textProvider.GenerateText(prompt, "gemini-2.0-flash")
}

// GenerateMarketingPlan خطة تسويق
func (s *StrategyService) GenerateMarketingPlan(product, budget, timeline string) (string, error) {
    prompt := fmt.Sprintf(`
    Create a comprehensive marketing plan for: %s
    
    Budget: %s
    Timeline: %s
    
    Include:
    1. Executive Summary
    2. Situation Analysis
    3. Marketing Objectives (SMART)
    4. Target Market
    5. Positioning Strategy
    6. Marketing Mix (4Ps)
    7. Action Plan
    8. Budget Breakdown
    9. Measurement & Evaluation
    
    Be practical and budget-conscious.
    `, product, budget, timeline)
    
    return s.textProvider.GenerateText(prompt, "llama3.2:3b")
}