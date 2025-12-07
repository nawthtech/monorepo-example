package services

import (
    "context"
    "fmt"
    "strings"
    "time"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type ContentService struct {
    textProvider types.TextProvider
}

func NewContentService(provider types.TextProvider) *ContentService {
    return &ContentService{
        textProvider: provider,
    }
}

func (s *ContentService) GenerateBlogPost(ctx context.Context, topic string, tone string, wordCount int) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Write a comprehensive blog post about: %s

Requirements:
- Tone: %s
- Word count: approximately %d words
- Include an engaging title
- Add relevant subheadings
- Include introduction, body, and conclusion
- Add a call-to-action at the end
- Optimize for SEO with relevant keywords
- Make it valuable and informative for readers

Structure:
1. Attention-grabbing title
2. Introduction (hook the reader)
3. Main content with subheadings
4. Practical tips or examples
5. Conclusion (summarize key points)
6. Call-to-action

Make it engaging and shareable.`, topic, tone, wordCount)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   wordCount * 2, // تقدير تقريبي
        Temperature: 0.8,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GenerateSocialMediaPost(ctx context.Context, platform string, message string, hashtags []string) (*types.TextResponse, error) {
    hashtagsStr := ""
    if len(hashtags) > 0 {
        hashtagsStr = "\nRelevant hashtags: " + strings.Join(hashtags, " ")
    }
    
    var platformGuidelines string
    switch strings.ToLower(platform) {
    case "twitter", "x":
        platformGuidelines = "Keep it under 280 characters. Use concise language."
    case "instagram":
        platformGuidelines = "Make it visual and engaging. Include emojis."
    case "facebook":
        platformGuidelines = "More detailed. Can be longer. Include questions to encourage comments."
    case "linkedin":
        platformGuidelines = "Professional tone. Focus on insights and value."
    case "tiktok":
        platformGuidelines = "Trendy and engaging. Use current slang if appropriate."
    default:
        platformGuidelines = "Keep it engaging and platform-appropriate."
    }
    
    prompt := fmt.Sprintf(`Write a social media post for %s:

Core message: %s
%s

Platform-specific requirements: %s

Include:
1. Main post content
2. Optional caption extension
3. Relevant emojis if appropriate
4. Call-to-action
5. %s

Make it engaging and platform-optimized.`, platform, message, hashtagsStr, platformGuidelines, hashtagsStr)
    
    maxTokens := 200
    if platform == "linkedin" || platform == "facebook" {
        maxTokens = 400
    }
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   maxTokens,
        Temperature: 0.9,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GenerateProductDescription(ctx context.Context, product string, features []string, targetCustomer string) (*types.TextResponse, error) {
    featuresStr := ""
    if len(features) > 0 {
        featuresStr = "Key features:\n"
        for _, feature := range features {
            featuresStr += fmt.Sprintf("- %s\n", feature)
        }
    }
    
    prompt := fmt.Sprintf(`Write a compelling product description for: %s

%s
Target customer: %s

Requirements:
1. Attention-grabbing headline
2. Clear value proposition
3. Highlight key benefits (not just features)
4. Use persuasive language
5. Include social proof elements
6. Clear call-to-action
7. Optimize for conversion

Tone: Professional yet conversational
Length: 150-300 words

Make it sell the benefits, not just the features.`, product, featuresStr, targetCustomer)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   500,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GenerateEmailNewsletter(ctx context.Context, topic string, audience string, length string) (*types.TextResponse, error) {
    var wordCount int
    switch length {
    case "short":
        wordCount = 300
    case "medium":
        wordCount = 600
    case "long":
        wordCount = 1000
    default:
        wordCount = 500
    }
    
    prompt := fmt.Sprintf(`Write an email newsletter about: %s

Target audience: %s
Length: approximately %d words

Structure:
1. Engaging subject line (5-10 words)
2. Pre-header text (complements subject line)
3. Opening greeting
4. Main content (value-driven)
5. Key takeaways
6. Call-to-action
7. Closing signature
8. P.S. section (optional but effective)

Requirements:
- Personal and conversational tone
- Mobile-friendly formatting
- Clear value proposition
- Include links where relevant
- Add personalization placeholders like {FirstName}

Make it valuable and click-worthy.`, topic, audience, wordCount)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   wordCount * 2,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GenerateAdCopy(ctx context.Context, product string, platform string, goal string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Write advertising copy for: %s

Platform: %s
Primary goal: %s

Include:
1. Headline (attention-grabbing)
2. Subheadline (supporting message)
3. Body copy (persuasive text)
4. Call-to-action (clear and compelling)
5. Display URL (if applicable)
6. Additional ad extensions ideas

Platform-specific requirements:
- Facebook/Instagram: Conversational, benefit-focused
- Google Ads: Keyword-rich, direct response
- LinkedIn: Professional, B2B focused
- Twitter: Concise, hashtag-friendly
- TikTok: Trendy, engaging, youth-focused

Make it conversion-optimized.`, product, platform, goal)
    
    maxTokens := 300
    if platform == "Twitter" || platform == "TikTok" {
        maxTokens = 150
    }
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   maxTokens,
        Temperature: 0.8,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GenerateVideoScript(ctx context.Context, videoType string, topic string, duration int) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Write a video script for: %s

Topic: %s
Target duration: %d seconds

Structure:
1. Hook (0-5 seconds: grab attention)
2. Introduction (5-15 seconds: set context)
3. Main content (15-45 seconds: deliver value)
4. Examples/demonstration (45-55 seconds: show, don't just tell)
5. Conclusion (55-60 seconds: summarize)
6. Call-to-action (last 5 seconds)

Include:
- Visual directions (what to show on screen)
- On-screen text suggestions
- Sound effects/music notes
- Pacing indicators
- Speaker notes/tonal guidance

Video type considerations:
- Explainer: Clear, educational
- Tutorial: Step-by-step, practical
- Promotional: Persuasive, benefit-focused
- Social media: Fast-paced, engaging

Make it visual and engaging.`, videoType, topic, duration)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   1000,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GenerateWhitepaperOutline(ctx context.Context, topic string, targetAudience string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Create a comprehensive whitepaper outline about: %s

Target audience: %s

Structure:
1. Title Page
   - Title
   - Subtitle
   - Author/Company
   - Date

2. Abstract/Executive Summary
   - Problem statement
   - Solution overview
   - Key findings

3. Table of Contents

4. Introduction
   - Background
   - Problem definition
   - Objectives
   - Methodology

5. Literature Review/Background
   - Current state
   - Related work
   - Market analysis

6. Main Analysis/Findings
   - Section 1: [Detailed analysis]
   - Section 2: [Data/results]
   - Section 3: [Case studies]
   - Section 4: [Technical details]

7. Discussion
   - Interpretation of findings
   - Implications
   - Limitations

8. Recommendations
   - Practical applications
   - Implementation steps
   - Future research

9. Conclusion
   - Summary
   - Final thoughts

10. References/Appendix
    - Citations
    - Additional data
    - Glossary

Make it thorough and research-oriented.`, topic, targetAudience)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   1500,
        Temperature: 0.6,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GeneratePressRelease(ctx context.Context, company string, announcement string, keyPoints []string) (*types.TextResponse, error) {
    keyPointsStr := ""
    if len(keyPoints) > 0 {
        keyPointsStr = "Key points to highlight:\n"
        for _, point := range keyPoints {
            keyPointsStr += fmt.Sprintf("- %s\n", point)
        }
    }
    
    prompt := fmt.Sprintf(`Write a professional press release for: %s

Announcement: %s

%s
Structure:
1. FOR IMMEDIATE RELEASE header
2. Headline (newsworthy and clear)
3. Dateline (CITY, STATE, Date)
4. Lead paragraph (who, what, when, where, why)
5. Body (details, quotes, context)
6. Boilerplate (company background)
7. Contact information
8. ### (end marker)

Requirements:
- Newsworthy angle
- Third-person perspective
- Include quotes from relevant executives
- Include relevant facts and figures
- Optimize for media pickup
- Include suggested social media posts

Make it journalistic and newsworthy.`, company, announcement, keyPointsStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   1000,
        Temperature: 0.6,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GenerateCaseStudy(ctx context.Context, client string, challenge string, solution string, results []string) (*types.TextResponse, error) {
    resultsStr := ""
    if len(results) > 0 {
        resultsStr = "Measurable results:\n"
        for _, result := range results {
            resultsStr += fmt.Sprintf("- %s\n", result)
        }
    }
    
    prompt := fmt.Sprintf(`Write a professional case study for: %s

Challenge/Problem: %s

Solution Provided: %s

%s
Structure:
1. Title (Client + Benefit)
2. Executive Summary
3. The Client (background)
4. The Challenge (detailed problem)
5. The Solution (implementation details)
6. The Results (quantifiable outcomes)
7. Testimonial/Quote
8. Conclusion
9. About [Your Company]

Requirements:
- Storytelling approach
- Data-driven results
- Client quotes (simulated if needed)
- Before/after comparison
- Visual elements suggestions
- ROI calculations

Make it compelling and results-focused.`, client, challenge, solution, resultsStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GenerateLandingPageCopy(ctx context.Context, product string, targetCustomer string, uniqueSellingPoints []string) (*types.TextResponse, error) {
    uspStr := ""
    if len(uniqueSellingPoints) > 0 {
        uspStr = "Unique selling points:\n"
        for _, usp := range uniqueSellingPoints {
            uspStr += fmt.Sprintf("- %s\n", usp)
        }
    }
    
    prompt := fmt.Sprintf(`Write high-converting landing page copy for: %s

Target customer: %s

%s
Sections to include:
1. Hero section (headline + subheadline + CTA)
2. Problem/Solution section
3. Benefits (not features)
4. Social proof/testimonials
5. Features/details
6. Pricing (if applicable)
7. FAQ section
8. Final CTA section

Requirements:
- Clear value proposition above the fold
- Benefit-oriented language
- Scannable formatting
- Strong calls-to-action
- Trust signals
- Mobile-optimized structure

Make it convert visitors into leads/customers.`, product, targetCustomer, uspStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *ContentService) GetServiceStats(ctx context.Context) map[string]interface{} {
    stats := make(map[string]interface{})
    
    // الحصول على إحصائيات من المزود إذا كانت متوفرة
    if provider, ok := s.textProvider.(interface{ GetStats() *types.ProviderStats }); ok {
        stats["provider_stats"] = provider.GetStats()
    }
    
    stats["service"] = "content"
    stats["content_types"] = []string{
        "blog_post",
        "social_media",
        "product_description", 
        "email_newsletter",
        "ad_copy",
        "video_script",
        "whitepaper",
        "press_release",
        "case_study",
        "landing_page",
    }
    stats["timestamp"] = time.Now()
    
    return stats
}